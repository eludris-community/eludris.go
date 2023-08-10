package types

import "errors"

var (
	ErrNoHttpUrl               = errors.New("HttpUrl is required")
	ErrGatewayAlreadyConnected = errors.New("gateway is already connected")
	ErrGatewayNotConnected     = errors.New("not connected ŧo the gateway")
	ErrNoGateway               = errors.New("no gateway configured")
)
