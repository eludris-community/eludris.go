// SPDX-License-Identifier: MIT

package events

import (
	"github.com/eludris-community/eludris.go/interfaces"
)

// EventListener represents a listener for an event.
type EventListener interface {
	Handle(event Event, client interfaces.Client)
}

type listenerFunc[E Event] struct {
	f func(E, interfaces.Client)
}

func (l listenerFunc[E]) Handle(event Event, client interfaces.Client) {
	// Try casting the event to the type we're expecting.
	// If it fails, this listener doesn't care about this event.
	if event, ok := event.(E); ok {
		l.f(event, client)
	}
}

func NewEventListenerFunc[E Event](f func(E, interfaces.Client)) EventListener {
	return &listenerFunc[E]{f: f}
}
