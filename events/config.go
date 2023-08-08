package events

import (
	"github.com/eludris-community/eludris.go/v2/interfaces"
)

type EventManagerConfig struct {
	Listeners []EventListener
}

func (c *EventManagerConfig) Apply(opts []EventManagerOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

type EventManagerOpt func(manager *EventManagerConfig)

func WithListenerFunc[E Event](f func(E, interfaces.Client)) EventManagerOpt {
	return WithEventListener[E](NewEventListenerFunc(f))
}

func WithEventListener[E Event](listener EventListener) EventManagerOpt {
	return func(manager *EventManagerConfig) {
		manager.Listeners = append(manager.Listeners, listener)
	}
}

func DefaultEventManagerConfig() *EventManagerConfig {
	return &EventManagerConfig{}
}
