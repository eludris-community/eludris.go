// SPDX-License-Identifier: MIT

// Package client provides a client for the Eludris API.
package client

import (
	"context"

	"github.com/eludris-community/eludris.go/v2/gateway"
	"github.com/eludris-community/eludris.go/v2/rest"
	"github.com/eludris-community/eludris.go/v2/types"
)

// Client represents a client for Eludris, with functions to interact with the API.
type Client interface {
	ConnectGateway(ctx context.Context) error
	EventManager() EventManager
	Rest() rest.Rest
}

type clientImpl struct {
	httpUrl      string
	gateway      gateway.Gateway
	rest         rest.Rest
	eventManager EventManager
}

// BuildClient builds a client from a configuration.
func BuildClient(config *Config) (Client, error) {
	if config == nil {
		return nil, types.ErrNoConfig
	}

	if config.Token == "" {
		return nil, types.ErrNoToken
	}

	if config.ApiUrl == "" {
		return nil, types.ErrNoApiUrl
	}

	client := &clientImpl{}

	if config.EventManager == nil {
		config.EventManager = NewEventManager(client, config.EventManagerOpts...)
	}
	client.eventManager = config.EventManager

	client.httpUrl = config.ApiUrl
	if config.WsUrl == "" || config.FileUrl == "" {
		tmp := rest.New(config.ApiUrl, "")
		info, err := tmp.GetInstanceInfo()
		if err != nil {
			return nil, err
		}

		if config.WsUrl == "" {
			config.WsUrl = info.PandemoniumUrl
		}
		if config.FileUrl == "" {
			config.FileUrl = info.EffisUrl
		}
	}
	if config.Rest == nil {
		config.Rest = rest.New(config.ApiUrl, config.FileUrl, config.RestOpts...)
	}
	client.rest = config.Rest

	if config.Gateway == nil && len(config.GatewayOpts) > 0 {
		config.GatewayOpts = append([]gateway.ConfigOpt{
			gateway.WithURL(config.WsUrl),
		}, config.GatewayOpts...)

		config.Gateway = gateway.New(config.Token, config.EventManager.HandleGatewayEvent, config.GatewayOpts...)
	}
	client.gateway = config.Gateway

	return client, nil
}

func (c *clientImpl) Rest() rest.Rest {
	return c.rest
}

func (c *clientImpl) EventManager() EventManager {
	return c.eventManager
}

func (c *clientImpl) ConnectGateway(ctx context.Context) error {
	if c.gateway == nil {
		return types.ErrNoGateway
	}

	return c.gateway.Connect(ctx)
}
