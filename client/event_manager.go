// SPDX-License-Identifier: MIT

package client

import (
	"github.com/apex/log"
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
)

type Event interface {
	Client() Client
}

// EventListener represents a listener for an event.
type EventListener interface {
	Handle(event Event)
}

type listenerFunc[E Event] struct {
	f func(E)
}

func (l listenerFunc[E]) Handle(event Event) {
	// Try casting the event to the type we're expecting.
	// If it fails, this listener doesn't care about this event.
	if event, ok := event.(E); ok {
		l.f(event)
	}
}

func NewEventListenerFunc[E Event](f func(E)) EventListener {
	return &listenerFunc[E]{f: f}
}

// EventManager manages events, dispatching and allowing subscriptions to events,
type EventManager interface {
	// Dispatch an event to all subscribers. Client is sent as the 2nd option to all subscribers.
	Dispatch(event Event)
	AddListeners(listeners ...EventListener)
	RemoveListeners(listeners ...EventListener)
	HandleGatewayEvent(gatewayEventType pandemonium.OpcodeType, event any)
}

type managerImpl struct {
	config EventManagerConfig
	client Client
}

func (m *managerImpl) AddListeners(listeners ...EventListener) {
	m.config.Listeners = append(m.config.Listeners, listeners...)
}

func (m *managerImpl) RemoveListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		for i, l := range m.config.Listeners {
			if l == listener {
				m.config.Listeners = append(m.config.Listeners[:i], m.config.Listeners[i+1:]...)
				break
			}
		}
	}
}

// Dispatch an event to all subscribers.
func (m *managerImpl) Dispatch(event Event) {
	for i := range m.config.Listeners {
		go func(i int) {
			m.config.Listeners[i].Handle(event)
		}(i)
	}
}

func (m *managerImpl) HandleGatewayEvent(eventType pandemonium.OpcodeType, event any) {
	if handler, ok := m.config.GatewayHandlers[eventType]; ok {
		handler.HandleGatewayEvent(m.client, event)
	} else {
		log.WithField("event_type", eventType).WithField("event", event).WithField("client", m.client).WithField("gateway_handlers", m.config.GatewayHandlers).WithField("gateway_handler", handler).WithField("ok", ok).WithField("event_type", eventType).WithField("event", event).WithField("client", m.client).WithField("gateway_handlers", m.config.GatewayHandlers).WithField("gateway_handler", handler).WithField("ok", ok).Error("handler for gateway event not found")
		log.Warnf("handler for gateway event '%s' not found", eventType)
	}
}

// NewEventManager creates a new event manager.
func NewEventManager(client Client, opts ...EventManagerOpt) EventManager {
	config := DefaultEventManagerConfig()
	config.Apply(opts)

	m := &managerImpl{client: client, config: *config}

	return m
}

type GatewayEventHandler interface {
	EventType() pandemonium.OpcodeType
	HandleGatewayEvent(client Client, event any)
}

func NewGatewayEventHandler[T any](eventType pandemonium.OpcodeType, handleFunc func(client Client, event T)) GatewayEventHandler {
	return &genericGatewayEventHandler[T]{eventType: eventType, handleFunc: handleFunc}
}

type genericGatewayEventHandler[T any] struct {
	eventType  pandemonium.OpcodeType
	handleFunc func(client Client, event T)
}

func (h *genericGatewayEventHandler[T]) EventType() pandemonium.OpcodeType {
	return h.eventType
}

func (h *genericGatewayEventHandler[T]) HandleGatewayEvent(client Client, event any) {
	if e, ok := event.(T); ok {
		h.handleFunc(client, e)
	}
}
