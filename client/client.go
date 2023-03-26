// SPDX-License-Identifier: MIT

// Package client provides a client for the Eludris API.
package client

import (
	"net/http"

	"github.com/eludris-community/eludris.go/events"
	"github.com/eludris-community/eludris.go/interfaces"
)

type clientImpl struct {
	httpUrl      string
	wsUrl        string
	fileUrl      string
	httpClient   http.Client
	eventManager events.EventManager
	rateLimiter  RateLimiter
}

// New creates a new client.
func New(config Config) interfaces.Client {
	if config.HTTPUrl == "" {
		config.HTTPUrl = "https://api.eludris.gay"
	}
	if config.WSUrl == "" {
		config.WSUrl = "wss://ws.eludris.gay/"
	}
	if config.FileUrl == "" {
		config.FileUrl = "https://cdn.eludris.gay"
	}
	if config.EventManager == nil {
		config.EventManager = events.NewEventManager()
	}
	if config.RateLimiter == nil {
		config.RateLimiter = NewRateLimiter()
	}

	return clientImpl{
		httpUrl:      config.HTTPUrl,
		wsUrl:        config.WSUrl,
		fileUrl:      config.FileUrl,
		httpClient:   http.Client{},
		eventManager: config.EventManager,
		rateLimiter:  config.RateLimiter,
	}
}
