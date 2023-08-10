package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/apex/log"
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
	"github.com/eludris-community/eludris.go/v2/types"
	"github.com/gorilla/websocket"
)

func New(eventHandlerFunc EventHandlerFunc, opts ...ConfigOpt) Gateway {
	config := DefaultConfig()
	config.Apply(opts)

	return &gatewayImpl{
		config:           *config,
		eventHandlerFunc: eventHandlerFunc,
	}
}

type gatewayImpl struct {
	config           Config
	eventHandlerFunc EventHandlerFunc

	conn      *websocket.Conn
	connMutex sync.Mutex

	lastPingSent     time.Time
	lastPongReceived time.Time
	pingInterval     time.Duration
	pingChan         chan struct{}
}

func (g *gatewayImpl) Connect(ctx context.Context) error {
	return g.connectTries(ctx, 0)
}

func (g *gatewayImpl) connectTries(ctx context.Context, try int) error {
	delay := time.Duration(try) * 2 * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	timer := time.NewTimer(delay)
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	if err := g.connect(ctx); err != nil {
		if errors.Is(err, types.ErrGatewayAlreadyConnected) {
			return err
		}
		log.WithError(err).Error("error connecting to gateway")
		return g.connectTries(ctx, try+1)
	}
	return nil
}

func (g *gatewayImpl) connect(ctx context.Context) error {
	log.WithField("url", g.config.URL).Info("Connecting to gateway")
	g.connMutex.Lock()
	defer g.connMutex.Unlock()

	if g.conn != nil {
		return types.ErrGatewayAlreadyConnected
	}

	url := g.config.URL
	g.lastPingSent = time.Now().UTC()
	conn, rs, err := g.config.Dialer.DialContext(ctx, url, nil)
	if err != nil {
		g.Close(ctx)
		body := "empty"
		if rs != nil && rs.Body != nil {
			defer func() {
				_ = rs.Body.Close()
			}()
			rawBody, bErr := io.ReadAll(rs.Body)
			if bErr != nil {
				log.WithField("body", body).Error("error while reading response body")
			}
			body = string(rawBody)
		}

		log.WithFields(
			log.Fields{
				"url":  url,
				"body": body,
			},
		).WithError(err).Error("error connecting to gateway")
		return err
	}

	g.conn = conn

	go g.listen(conn)

	return nil
}

func (g *gatewayImpl) reconnect() {
	err := g.connectTries(context.Background(), 0)
	if err != nil {
		log.WithError(err).Error("error reconnecting to gateway")
	}
}

func (g *gatewayImpl) listen(conn *websocket.Conn) {
loop:
	for {
		mt, data, err := conn.ReadMessage()
		if err != nil {
			g.connMutex.Lock()
			sameConnection := g.conn == conn
			g.connMutex.Unlock()

			if !sameConnection {
				return
			}

			reconnect := true
			var closeError *websocket.CloseError
			if errors.Is(err, net.ErrClosed) {
				reconnect = false
			} else if !errors.As(err, &closeError) {
				log.WithError(err).Debug("failed to read next message from gateway")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			g.CloseWithCode(ctx, websocket.CloseServiceRestart, "reconnecting")
			cancel()
			if reconnect {
				go g.reconnect()
			}

			break loop
		}

		message, err := g.parseMessage(mt, data)
		if err != nil {
			log.WithError(err).Error("error while parsing gateway message")
			continue
		}

		if message.Op == pandemonium.HelloOp {
			g.pingInterval = time.Duration(message.D.(pandemonium.Hello).HeartbeatInterval) * time.Millisecond
			g.lastPongReceived = time.Now().UTC()
			go g.ping()
		}

		g.eventHandlerFunc(message.Op, message.D)
	}
}

func (g *gatewayImpl) ping() {
	if g.pingChan == nil {
		g.pingChan = make(chan struct{})
	}
	pingTicker := time.NewTicker(g.pingInterval)
	defer pingTicker.Stop()
	defer log.Debug("stopping pinging")

	for {
		select {
		case <-pingTicker.C:
			go g.doPing()
		case <-g.pingChan:
			return
		}
	}
}

func (g *gatewayImpl) doPing() {
	log.Debug("sending ping")

	ctx, cancel := context.WithTimeout(context.Background(), g.pingInterval)
	defer cancel()
	if err := g.Send(ctx, pandemonium.PingOp, nil); err != nil {
		if errors.Is(err, types.ErrGatewayNotConnected) || errors.Is(err, syscall.EPIPE) {
			return
		}
		log.WithError(err).Error("error sending ping")
		g.CloseWithCode(context.TODO(), websocket.CloseServiceRestart, "ping timeout")
		go g.reconnect()
		return
	}
	g.lastPingSent = time.Now().UTC()
}

func (g *gatewayImpl) Send(ctx context.Context, op pandemonium.OpcodeType, d any) error {
	data, err := json.Marshal(Payload{
		Op: op,
		D:  d,
	})
	if err != nil {
		return err
	}
	return g.send(ctx, websocket.TextMessage, data)
}

func (g *gatewayImpl) send(ctx context.Context, messageType int, data []byte) error {
	g.connMutex.Lock()
	defer g.connMutex.Unlock()
	if g.conn == nil {
		return types.ErrGatewayNotConnected
	}

	log.WithField("message", string(data)).Debug("sending gateway message")
	return g.conn.WriteMessage(messageType, data)
}

func (g *gatewayImpl) Latency() time.Duration {
	return g.lastPongReceived.Sub(g.lastPingSent)
}

func (g *gatewayImpl) parseMessage(mt int, data []byte) (Payload, error) {
	log.WithField("message", string(data)).Debug("received gateway message")

	var payload Payload
	return payload, json.Unmarshal(data, &payload)
}

func (g *gatewayImpl) CloseWithCode(ctx context.Context, code int, message string) {
	g.connMutex.Lock()
	defer g.connMutex.Unlock()
	if g.conn != nil {
		log.WithFields(log.Fields{
			"code":    code,
			"message": message,
		}).Info("Closing gateway connection")
		if err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message)); err != nil {
			log.WithError(err).Error("error closing gateway connection")
		}
		_ = g.conn.Close()
		g.conn = nil
	}
}

func (g *gatewayImpl) Close(ctx context.Context) {
	g.CloseWithCode(ctx, websocket.CloseNormalClosure, "Shutting down")
}
