package types

import "errors"

var (
	ErrNoApiUrl                = errors.New("ApiUrl is required")
	ErrGatewayAlreadyConnected = errors.New("gateway is already connected")
	ErrGatewayNotConnected     = errors.New("not connected ŧo the gateway")
	ErrNoGateway               = errors.New("no gateway configured")
)
