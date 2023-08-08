// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"os/signal" // whatsapp
	"syscall"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/eludris-community/eludris.go/v2/client"
	"github.com/eludris-community/eludris.go/v2/events"
	"github.com/eludris-community/eludris.go/v2/interfaces"
)

func onMessage(msg *events.MessageEvent, c interfaces.Client) {
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

	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.Default)

	manager := events.NewEventManager()
	events.Subscribe(manager, onMessage)
	c, err := client.New(
		client.WithEventManagerOpts(events.WithListenerFunc(onMessage)),
		client.WithHttpUrl(HTTPUrl),
	)

	if err != nil {
		panic(err)
	}

	if err := c.Connect(); err != nil {
		panic(err)
	}

	// // Ratelimit test
	// for i := 0; i < 20; i++ {
	// 	go func(i int) {
	// 		msg, err := c.SendMessage("hewwo from gowang", fmt.Sprintf("hello %d", i))
	// 		if err != nil {
	// 			fmt.Printf("Error: %s", err.Error())
	// 		}
	// 		fmt.Println(msg)
	// 	}(i)
	// }

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
