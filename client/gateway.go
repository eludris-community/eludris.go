// SPDX-License-Identifier: MIT

package client

import (
	"time"

	"github.com/apex/log"
	"github.com/eludris-community/eludris.go/events"
	"github.com/eludris-community/eludris.go/interfaces"
	"github.com/gorilla/websocket"
)

// Connect connects to the Eludris websocket for events.
func (c clientImpl) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.wsUrl, nil)

	if err != nil {
		return err
	}

	done := make(chan struct{})
	pongs := make(chan struct{})

	dispatchPong := func(*events.PongEvent, interfaces.Client) {
		pongs <- struct{}{}
	}

	go func() {
		defer close(done)
		events.Subscribe(c.eventManager, dispatchPong)
		go ping(conn, pongs)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.WithError(err).Error("Error reading message")
				err := c.Connect()
				if err != nil {
					panic(err)
				}

				return
			}
			log.WithField("message", string(message)).Debug("Received message")
			c.eventManager.Dispatch(c, message)
		}
	}()

	return nil
}

func ping(conn *websocket.Conn, pongs chan struct{}) {
	for {
		time.Sleep(44 * time.Second)
		conn.WriteMessage(websocket.TextMessage, []byte("{\"op\": \"PING\"}"))

		select {

		case <-pongs:
			continue
		case <-time.After(1 * time.Second):
			conn.Close()
		}
	}
}
