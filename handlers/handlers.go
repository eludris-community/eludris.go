package handlers

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
	"github.com/eludris-community/eludris.go/v2/client"
)

func AllHandlers() map[pandemonium.OpcodeType]client.GatewayEventHandler {
	handlers := make(map[pandemonium.OpcodeType]client.GatewayEventHandler, len(allEventHandlers))
	for _, handler := range allEventHandlers {
		handlers[handler.EventType()] = handler
	}
	return handlers
}

var allEventHandlers = []client.GatewayEventHandler{
	client.NewGatewayEventHandler(pandemonium.PongOp, gatewayHandlerPong),
	client.NewGatewayEventHandler(pandemonium.HelloOp, gatewayHandlerHello),
	client.NewGatewayEventHandler(pandemonium.RatelimitOp, gatewayHandlerRateLimit),
	client.NewGatewayEventHandler(pandemonium.MessageCreateOp, gatewayHandlerMessageCreate),
}
