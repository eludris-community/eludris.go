package events

import (
	"github.com/eludris-community/eludris.go/interfaces"
)

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

type MessageEvent struct {
	event
	Content string `mapstructure:"content"`
	Author  string `mapstructure:"author"`
}

func (*MessageEvent) Op() string {
	return "MESSAGE_CREATE"
}

// Pong has no inner data.
type PongEvent struct {
	event
}

func (*PongEvent) Op() string {
	return "PONG"
}
