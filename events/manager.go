// SPDX-License-Identifier: MIT

package events

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/apex/log"
	"github.com/eludris-community/eludris.go/interfaces"
	"github.com/mitchellh/mapstructure"
)

// EventListener represents a listener for an event.
type EventListener interface {
	// Handle handles an event.
	Handle(Event, interfaces.Client)
	// Op returns the opcode for this listener.
	Op() string
	// Func returns the inner function for this limiter.
	Func() func(Event, interfaces.Client)
}

// eventListener is the internal implementation of EventListener
type eventListener[E Event] struct {
	op string
	f  func(E, interfaces.Client)
}

func (e eventListener[E]) Op() string {
	return e.op
}

func (e eventListener[E]) Func() func(Event, interfaces.Client) {
	// Turn the inner function into a normal `Event` from `E`.
	return any(e.f).(func(Event, interfaces.Client))
}

func (l eventListener[E]) Handle(event Event, client interfaces.Client) {
	if event, ok := event.(E); ok {
		l.f(event, client)
	}
}

// EventManager manages events, dispatching and allowing subscriptions to events,
type EventManager interface {
	// Subscribe to an event with an event listener. The event is taken via reflection.
	Subscribe(EventListener)
	// Dispatch an event to all subscribers. Client is sent as the 2nd option to all subscribers.
	Dispatch(interfaces.Client, []byte)
}

type managerImpl struct {
	subscribers map[string][]EventListener
}

// Subscribe allows you to subscribe to an event to the given manager.
// This infers the event type from the function signature.
func Subscribe[E Event](m EventManager, subscriber func(E, interfaces.Client)) {
	t := reflect.TypeOf(subscriber)
	eventT := t.In(0)

	// Get the actual Op() method.
	for i := 0; i < eventT.NumMethod(); i++ {
		method := eventT.Method(i)
		if method.Name == "Op" {
			op := method.Func.Call([]reflect.Value{reflect.New(eventT).Elem()})[0].String()
			listener := eventListener[E]{op: op, f: subscriber}
			m.Subscribe(&listener)
			return
		}
	}
}

func (m *managerImpl) Subscribe(listener EventListener) {
	// Create the slice if it doesn't exist.
	if _, ok := m.subscribers[listener.Op()]; !ok {
		m.subscribers[listener.Op()] = make([]EventListener, 0)
	}

	m.subscribers[listener.Op()] = append(m.subscribers[listener.Op()], listener)
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

	subscribers := m.subscribers[op]

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

	for _, subscriber := range subscribers {
		subscriber.Handle(event, client)
	}
}

// NewEventManager creates a new event manager.
func NewEventManager() EventManager {
	return &managerImpl{subscribers: make(map[string][]EventListener)}
}
