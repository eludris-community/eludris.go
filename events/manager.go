// SPDX-License-Identifier: MIT

package events

import (
	"bytes"
	"encoding/json"

	"github.com/apex/log"
	"github.com/eludris-community/eludris.go/v2/interfaces"
	"github.com/mitchellh/mapstructure"
)

// EventManager manages events, dispatching and allowing subscriptions to events,
type EventManager interface {
	// Subscribe to an event with an event listener. The event is taken via reflection.
	Subscribe(EventListener)
	// Dispatch an event to all subscribers. Client is sent as the 2nd option to all subscribers.
	Dispatch(interfaces.Client, []byte)
}

type managerImpl struct {
	config EventManagerConfig
}

// Subscribe allows you to subscribe to an event to the given manager.
// This infers the event type from the function signature.
func Subscribe[E Event](m EventManager, subscriber func(E, interfaces.Client)) {
	listener := NewEventListenerFunc(subscriber)

	m.Subscribe(listener)
}

func (m *managerImpl) Subscribe(listener EventListener) {
	m.config.Listeners = append(m.config.Listeners, listener)
}

// Dispatch an event to all subscribers.
func (m *managerImpl) Dispatch(client interfaces.Client, data []byte) {
	// Treat it as a map at first.
	var msg map[string]any
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&msg)
	log.WithField("msg", msg).Debug("Dispatching event")

	op := msg["op"].(string)
	innerData := msg["d"]
	if innerData == nil {
		// This is the actual data we want to decode.
		innerData = make(map[string]any)
	}

	var event Event

	switch op {

	case "MESSAGE_CREATE":
		{
			var innerEvent MessageEvent
			mapstructure.Decode(innerData, &innerEvent)
			event = &innerEvent
		}
	case "PONG":
		{
			var innerEvent PongEvent
			mapstructure.Decode(innerData, &innerEvent)
			event = &innerEvent
		}
	default:
		log.WithField("op", op).Warn("Unknown event type")

		return
	}

	for i := range m.config.Listeners {
		go func(i int) {
			m.config.Listeners[i].Handle(event, client)
		}(i)
	}
}

// NewEventManager creates a new event manager.
func NewEventManager(opts ...EventManagerOpt) EventManager {
	config := DefaultEventManagerConfig()
	config.Apply(opts)

	m := &managerImpl{config: *config}

	return m
}
