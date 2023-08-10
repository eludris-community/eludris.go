package handlers

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
	"github.com/eludris-community/eludris.go/v2/client"
	"github.com/eludris-community/eludris.go/v2/events"
)

func gatewayHandlerMessageCreate(client client.Client, event pandemonium.MessageCreate) {
	client.EventManager().Dispatch(&events.MessageCreate{GenericEvent: events.NewGenericEvent(client), MessageCreate: event})
}
