// SPDX-License-Identifier: MIT

package events

import (
	"github.com/eludris-community/eludris.go/interfaces"
)

// Event represents the base fields for all events.
type Event interface {
	Op() string
	Client() interfaces.Client
	SetClient(interfaces.Client)
}

type event struct {
	client interfaces.Client
}

func (e event) Client() interfaces.Client {
	return e.client
}

func (e *event) SetClient(client interfaces.Client) {
	e.client = client
}

// MessageEvent represents a received message from Eludris.
type MessageEvent struct {
	event
	Content string `mapstructure:"content"`
	Author  string `mapstructure:"author"`
}

func (*MessageEvent) Op() string {
	return "MESSAGE_CREATE"
}

type PongEvent struct {
	event
}

func (*PongEvent) Op() string {
	return "PONG"
}
