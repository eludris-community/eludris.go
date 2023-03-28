// SPDX-License-Identifier: MIT

package client

import "github.com/eludris-community/eludris.go/events"

// Config represents configuration for the client.
// This allows you to set the URLs for your chosen Eludris instance.
// These contain defaults, notably to the "official" Eludris instance https://eludris.gay.
type Config struct {
	// HTTPUrl is the URL for the HTTP API, A.K.A "oprish".
	HTTPUrl string
	// WSUrl is the URL for the WebSocket gateway, A.K.A "pandemonium".
	WSUrl string
	// FileUrl is the URL for the file server, A.K.A "effis".
	FileUrl string
	// An event manager to use for the client.
	EventManager events.EventManager
	// A rate limiter to use for the client.
	RateLimiter RateLimiter
}
