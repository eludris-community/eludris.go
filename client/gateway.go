// SPDX-License-Identifier: MIT

package client

import (
	"fmt"
	"time"

	"github.com/eludris-community/eludris.go/events"
	"github.com/gorilla/websocket"
)

func (c clientImpl) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.wsUrl, nil)

	if err != nil {
		return err
	}

	done := make(chan struct{})
	pongs := make(chan struct{})

	dispatchPing := func(*events.PongEvent) {
		pongs <- struct{}{}
	}

	go func() {
		defer close(done)
		events.Subscribe(c.eventManager, dispatchPing)
		go ping(conn, pongs)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("error: %v\n", err)
				err := c.Connect()
				if err != nil {
					panic(err)
				}

				return
			}
			fmt.Printf("recv: %s\n", message)
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
