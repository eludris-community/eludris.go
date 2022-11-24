package events

import (
	"github.com/ooliver1/eludris.go/interfaces"
)

type Event struct {
	client interfaces.Client
}

func (e *Event) Client() interfaces.Client {
	return e.client
}

func (e *Event) SetClient(client interfaces.Client) {
	e.client = client
}

type MessageEvent struct {
	Event
	Content string `json:"content"`
	Author  string `json:"author"`
}
