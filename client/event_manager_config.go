package client

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
)

type EventManagerConfig struct {
	Listeners       []EventListener
	GatewayHandlers map[pandemonium.OpcodeType]GatewayEventHandler
}

func (c *EventManagerConfig) Apply(opts []EventManagerOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

type EventManagerOpt func(config *EventManagerConfig)

func WithListenerFunc[E Event](f func(E)) EventManagerOpt {
	return WithEventListener[E](NewEventListenerFunc(f))
}

func WithEventListener[E Event](listener EventListener) EventManagerOpt {
	return func(manager *EventManagerConfig) {
		manager.Listeners = append(manager.Listeners, listener)
	}
}

func WithGatewayHandlers(handlers map[pandemonium.OpcodeType]GatewayEventHandler) EventManagerOpt {
	return func(config *EventManagerConfig) {
		config.GatewayHandlers = handlers
	}
}

func DefaultEventManagerConfig() *EventManagerConfig {
	return &EventManagerConfig{}
}
