package main

import (
	// "fmt"
	// "time"
	"os"
	"os/signal" // whatsapp
	"syscall"

	"github.com/ooliver1/eludris.go/client"
	"github.com/ooliver1/eludris.go/events"
)

func onMessage(msg events.MessageEvent) {
	c := msg.Client()

	if msg.Author == "hello" {
		return
	}

	_, err := c.SendMessage("hewwo from gowang", "hello")
	if err != nil {
		panic(err)
	}
}

func main() {
	manager := events.NewEventManager()
	manager.Subscribe(onMessage)
	c := client.New(client.Config{EventManager: manager})
	err := c.Connect()

	if err != nil {
		panic(err)
	}

	// for {
	// 	msg, err := c.SendMessage("hewwo from gowang", "hello")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(msg)
	// 	time.Sleep(1 * time.Second)
	// }

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
