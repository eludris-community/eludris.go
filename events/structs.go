// SPDX-License-Identifier: MIT

package events

// Event represents the base fields for all events.
type Event interface {
	Op() string
}

// MessageEvent represents a received message from Eludris.
type MessageEvent struct {
	Content string `mapstructure:"content"`
	Author  string `mapstructure:"author"`
}

func (*MessageEvent) Op() string {
	return "MESSAGE_CREATE"
}

type PongEvent struct{}

func (*PongEvent) Op() string {
	return "PONG"
}
