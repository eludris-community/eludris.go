package handlers

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
	"github.com/eludris-community/eludris.go/v2/client"
	"github.com/eludris-community/eludris.go/v2/events"
)

func gatewayHandlerUserCreate(client client.Client, event pandemonium.MessageCreate) {
	client.EventManager().Dispatch(&events.MessageCreate{GenericEvent: events.NewGenericEvent(client), MessageCreate: event})
}

func gatewayHandlerUserUpdate(client client.Client, event pandemonium.UserUpdate) {
	client.EventManager().Dispatch(&events.UserUpdate{GenericEvent: events.NewGenericEvent(client), UserUpdate: event})
}

func gatewayHandlerPresenceUpdate(client client.Client, event pandemonium.PresenceUpdate) {
	client.EventManager().Dispatch(&events.PresenceUpdate{GenericEvent: events.NewGenericEvent(client), PresenceUpdate: event})
}
