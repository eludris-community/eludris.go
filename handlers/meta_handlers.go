package handlers

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
	"github.com/eludris-community/eludris.go/v2/client"
	"github.com/eludris-community/eludris.go/v2/events"
)

func gatewayHandlerPong(client client.Client, event interface{}) {
	client.EventManager().Dispatch(&events.Pong{GenericEvent: events.NewGenericEvent(client)})
}

func gatewayHandlerHello(client client.Client, event pandemonium.Hello) {
	client.EventManager().Dispatch(&events.Hello{GenericEvent: events.NewGenericEvent(client), Hello: event})
}

func gatewayHandlerRateLimit(client client.Client, event pandemonium.Ratelimit) {
	client.EventManager().Dispatch(&events.RateLimit{GenericEvent: events.NewGenericEvent(client), Ratelimit: event})
}

func gatewayHandlerAuthenticated(client client.Client, event pandemonium.Authenticated) {
	client.EventManager().Dispatch(&events.Authenticated{GenericEvent: events.NewGenericEvent(client), Authenticated: event})
}
