package client

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func (c clientImpl) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.wsUrl, nil)

	if err != nil {
		return err
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		go ping(conn)
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			fmt.Printf("recv: %s", message)
			println(string(message))
			fmt.Printf("dispatching event: %s\n", message)
			c.eventManager.Dispatch(c, message)
		}
	}()

	return nil
}

func ping(conn *websocket.Conn) {
	for {
		time.Sleep(20 * time.Second)
		conn.WriteMessage(websocket.PingMessage, []byte{})
	}
}
