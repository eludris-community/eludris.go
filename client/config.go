package client

import "github.com/ooliver1/eludris.go/events"

type Config struct {
	HTTPUrl      string
	WSUrl        string
	FileUrl      string
	EventManager events.EventManager
}
