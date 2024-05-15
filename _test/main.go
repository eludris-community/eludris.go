// SPDX-License-Identifier: MIT

package main

import (
	"context"
	// "fmt"
	"os"
	"os/signal" // whatsapp
	"syscall"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/eludris-community/eludris.go/v2"
	"github.com/eludris-community/eludris.go/v2/client"
	// "github.com/eludris-community/eludris.go/v2/events"
	// "github.com/eludris-community/eludris.go/v2/types"

	"github.com/joho/godotenv"
)

// func onMessage(msg *events.MessageCreate) {
// 	c := msg.Client()
// 	if msg.Author == "hello" {
// 		return
// 	}

// 	switch msg.Content {
// 	case "!speed":
// 		_, err := c.Rest().CreateMessage(types.MessageCreate{Author: "googwuki", Content: "I am the fastest ever."})

// 		if err != nil {
// 			panic(err)
// 		}
// 	case "!cat":
// 		file, err := os.Open("cat.jpg")
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer file.Close()

// 		data, err := c.Rest().UploadAttachment(file, false)
// 		if err != nil {
// 			panic(err)
// 		}
// 		log.Infof("%v", data)
// 		_, err = c.Rest().CreateMessage(types.MessageCreate{Author: "googwuki", Content: fmt.Sprintf("https://cdn.eludris.gay/%s", data.Id)})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.Default)

	c, err := eludris.New(
		os.Getenv("ELUDRIS_TOKEN"),
		// client.WithEventManagerOpts(client.WithListenerFunc(onMessage)),
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
