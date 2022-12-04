// SPDX-License-Identifier: MIT

package client

import "github.com/eludris-community/eludris.go/events"

// Config represents configuration for the client.
// This allows you to set the URLs for your chosen Eludris instance.
type Config struct {
	HTTPUrl      string
	WSUrl        string
	FileUrl      string
	EventManager events.EventManager
}
