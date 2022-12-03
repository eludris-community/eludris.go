package main

import (
	"fmt"
	// "time"
	"os"
	"os/signal" // whatsapp
	"syscall"

	"github.com/eludris-community/eludris.go/client"
	"github.com/eludris-community/eludris.go/events"
)

func onMessage(msg *events.MessageEvent) {
	c := msg.Client()

	if msg.Author == "hello" {
		return
	}

	switch msg.Content {
	case "!speed":
		_, err := c.SendMessage("googwuki", "I am the fastest ever.")

		if err != nil {
			panic(err)
		}
	case "!cat":
		file, err := os.Open("cat.jpg")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		data, err := c.UploadAttachment(file, false)
		if err != nil {
			panic(err)
		}

		_, err = c.SendMessage("googwuki", fmt.Sprintf("https://cdn.eludris.gay/%s", data.Id))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	HTTPUrl := os.Getenv("ELUDRIS_HTTP_URL")
	WSUrl := os.Getenv("ELUDRIS_WS_URL")

	manager := events.NewEventManager()
	events.Subscribe(manager, onMessage)
	c := client.New(client.Config{HTTPUrl: HTTPUrl, WSUrl: WSUrl, EventManager: manager})
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
