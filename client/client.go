// SPDX-License-Identifier: MIT

// Package client provides a client for the Eludris API.
package client

import (
	"errors"
	"net/http"

	"github.com/eludris-community/eludris-api-types.go/models"
	"github.com/eludris-community/eludris.go/v2/events"
	"github.com/eludris-community/eludris.go/v2/interfaces"
)

var (
	ErrNoHttpUrl = errors.New("HttpUrl is required")
)

type clientImpl struct {
	httpUrl      string
	wsUrl        string
	fileUrl      string
	httpClient   http.Client
	eventManager events.EventManager
	rateLimiter  RateLimiter
}

// BuildClient builds a client from a configuration.
func BuildClient(config *Config) (interfaces.Client, error) {
	if config.HttpUrl == "" {
		return nil, ErrNoHttpUrl
	}

	if config.WsUrl == "" || config.FileUrl == "" {
		tmp := clientImpl{httpUrl: config.HttpUrl, rateLimiter: NewRateLimiter()}
		info, err := tmp.FetchInstanceInfo()
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

	if config.EventManager == nil {
		config.EventManager = events.NewEventManager(config.EventManagerOpts...)
	}

	if config.RateLimiter == nil {
		config.RateLimiter = NewRateLimiter(config.RateLimiterOpts...)
	}

	return &clientImpl{
		httpUrl:      config.HttpUrl,
		wsUrl:        config.WsUrl,
		fileUrl:      config.FileUrl,
		httpClient:   http.Client{},
		eventManager: config.EventManager,
		rateLimiter:  config.RateLimiter,
	}, nil
}

// New creates a new client.
func New(opts ...ConfigOpt) (interfaces.Client, error) {
	config := DefaultConfig()
	config.Apply(opts)

	return BuildClient(config)
}

// FetchInstanceInfo fetches information about the instance.
func (c *clientImpl) FetchInstanceInfo() (models.InstanceInfo, error) {
	var res models.InstanceInfo
	_, err := c.request(InstanceInfo.Compile(nil), Data{}, &res)

	return res, err
}
