package gateway

import (
	"context"
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
	"time"
)

// Gateway connects to the Eludris websocket gateway (pandemonium).
type Gateway interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context)
	CloseWithCode(ctx context.Context, code int, message string)
	Send(ctx context.Context, op pandemonium.OpcodeType, data any) error
	Latency() time.Duration
}

type EventHandlerFunc func(gatewayEventType pandemonium.OpcodeType, event any)
