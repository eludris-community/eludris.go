package client

import (
	"net/http"

	"github.com/ooliver1/eludris.go/events"
	"github.com/ooliver1/eludris.go/interfaces"
)

type clientImpl struct {
	httpUrl      string
	wsUrl        string
	httpClient   http.Client
	eventManager events.EventManager
}

func New(config Config) interfaces.Client {
	if config.HTTPUrl == "" {
		config.HTTPUrl = "https://api.eludris.gay"
	}
	if config.WSUrl == "" {
		config.WSUrl = "wss://ws.eludris.gay/"
	}
	if config.EventManager == nil {
		config.EventManager = events.NewEventManager()
	}

	return clientImpl{
		httpUrl:      config.HTTPUrl,
		wsUrl:        config.WSUrl,
		httpClient:   http.Client{},
		eventManager: config.EventManager,
	}
}
