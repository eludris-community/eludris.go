// SPDX-License-Identifier: MIT

package events

// Event represents the base implementation for all events.
type Event interface {
	Op() string
}

// MessageEvent represents a MESSAGE_CREATE event from Eludris.
type MessageEvent struct {
	Content string `mapstructure:"content"`
	Author  string `mapstructure:"author"`
}

func (*MessageEvent) Op() string {
	return "MESSAGE_CREATE"
}

// PongEvent represents a PONG event from Eludris.
type PongEvent struct{}

func (*PongEvent) Op() string {
	return "PONG"
}
