// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal" // whatsapp
	"syscall"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/eludris-community/eludris.go/v2"
	"github.com/eludris-community/eludris.go/v2/client"
	"github.com/eludris-community/eludris.go/v2/events"
)

func onMessage(msg *events.MessageCreate) {
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

	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.Default)

	c, err := eludris.New(
		client.WithEventManagerOpts(client.WithListenerFunc(onMessage)),
		client.WithHttpUrl(HTTPUrl),
		client.WithDefaultGateway(),
	)

	if err != nil {
		panic(err)
	}

	if err := c.ConnectGateway(context.TODO()); err != nil {
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
