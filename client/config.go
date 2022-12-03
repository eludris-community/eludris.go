package client

import "github.com/eludris-community/eludris.go/events"

type Config struct {
	HTTPUrl      string
	WSUrl        string
	FileUrl      string
	EventManager events.EventManager
}
