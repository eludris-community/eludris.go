// SPDX-License-Identifier: MIT

package client

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
	"github.com/eludris-community/eludris.go/v2/gateway"
	"github.com/eludris-community/eludris.go/v2/rest"
)

type ConfigOpt func(config *Config)

// Config represents configuration for the client.
// This allows you to set the URLs for your chosen Eludris instance.
// These contain defaults, notably to the "official" Eludris instance https://eludris.gay.
type Config struct {
	// Token is the session token to use in both the HTTP API and WebSocket gateway.
	Token string
	// ApiUrl is the URL for the HTTP API, A.K.A "oprish".
	ApiUrl string
	// WSUrl is the URL for the WebSocket gateway, A.K.A "pandemonium".
	// This is obtainable from HTTPUrl if not set.
	WsUrl string
	// FileUrl is the URL for the file server, A.K.A "effis".
	// This is obtainable from HTTPUrl if not set.
	FileUrl string
	// An event manager to use for the client.
	EventManager EventManager
	// Options to apply to EventManager
	EventManagerOpts []EventManagerOpt
	Gateway          gateway.Gateway
	GatewayOpts      []gateway.ConfigOpt
	Rest             rest.Rest
	RestOpts         []rest.ConfigOpt
}

// DefaultConfig returns a new Config with default values.
func DefaultConfig(gatewayHandlers map[pandemonium.OpcodeType]GatewayEventHandler) *Config {
	return &Config{
		EventManagerOpts: []EventManagerOpt{WithGatewayHandlers(gatewayHandlers)},
	}
}

// Apply applies the given ConfigOpts to the Config.
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithHttpUrl returns a ConfigOpt that sets the HTTP URL.
func WithApiUrl(url string) ConfigOpt {
	return func(config *Config) {
		config.ApiUrl = url
	}
}

// WithWsUrl returns a ConfigOpt that sets the WebSocket URL.
func WithWsUrl(url string) ConfigOpt {
	return func(config *Config) {
		config.WsUrl = url
	}
}

// WithFileUrl returns a ConfigOpt that sets the file URL.
func WithFileUrl(url string) ConfigOpt {
	return func(config *Config) {
		config.FileUrl = url
	}
}

// WithEventManager returns a ConfigOpt that sets the event manager.
func WithEventManager(manager EventManager) ConfigOpt {
	return func(config *Config) {
		config.EventManager = manager
	}
}

// WithEventManagerOpts returns a ConfigOpt that sets the event manager options.
func WithEventManagerOpts(opts ...EventManagerOpt) ConfigOpt {
	return func(config *Config) {
		config.EventManagerOpts = append(config.EventManagerOpts, opts...)
	}
}

// WithGateway lets you inject your own gateway.Gateway.
func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

func WithDefaultGateway() ConfigOpt {
	return func(config *Config) {
		config.GatewayOpts = append(config.GatewayOpts, func(_ *gateway.Config) {})
	}
}

func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayOpts = append(config.GatewayOpts, opts...)
	}
}
