// SPDX-License-Identifier: MIT

// Package client provides a client for the Eludris API.
package client

import (
	"context"
	"io"
	"net/http"

	"github.com/eludris-community/eludris-api-types.go/v2/models"
	"github.com/eludris-community/eludris.go/v2/gateway"
	"github.com/eludris-community/eludris.go/v2/types"
)

// Client represents a client for Eludris, with functions to interact with the API.
type Client interface {
	ConnectGateway(ctx context.Context) error
	EventManager() EventManager

	// SendMessage sends a message to the chat.
	// Returns the message that was sent.
	// Author must be between 1 and 32 characters long.
	// Content must be between 1 and the instance's configured max message length.
	SendMessage(author, content string) (models.Message, error)

	// UploadAttachment uploads a file to the instance's attachments bucket.
	UploadAttachment(file io.Reader, spoiler bool) (models.FileData, error)
	// UploadFile uploads a file to a specific bucket.
	// Currently only "attachments" exists.
	UploadFile(bucket string, file io.Reader, spoiler bool) (models.FileData, error)
	// FetchAttachment fetches a file from the instance's attachments bucket.
	FetchAttachment(id string) (io.ReadCloser, error)
	// FetchFile fetches a file from a specific bucket.
	// Currently only "attachments" exists.
	FetchFile(bucket, id string) (io.ReadCloser, error)
	// FetchAttachmentData fetches a file from the instance's attachments bucket.
	FetchAttachmentData(id string) (models.FileData, error)
	// FetchFileData fetches a file from a specific bucket.
	// Currently only "attachments" exists.
	FetchFileData(bucket, id string) (models.FileData, error)
}

type clientImpl struct {
	httpUrl      string
	fileUrl      string
	httpClient   http.Client
	gateway      gateway.Gateway
	eventManager EventManager
	rateLimiter  RateLimiter
}

// BuildClient builds a client from a configuration.
func BuildClient(config *Config) (Client, error) {
	if config.HttpUrl == "" {
		return nil, types.ErrNoHttpUrl
	}

	client := &clientImpl{}

	if config.EventManager == nil {
		config.EventManager = NewEventManager(client, config.EventManagerOpts...)
	}
	client.eventManager = config.EventManager

	if config.RateLimiter == nil {
		config.RateLimiter = NewRateLimiter(config.RateLimiterOpts...)
	}
	client.rateLimiter = config.RateLimiter

	client.httpUrl = config.HttpUrl
	if config.WsUrl == "" || config.FileUrl == "" {
		info, err := client.FetchInstanceInfo()
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

	if config.Gateway == nil && len(config.GatewayOpts) > 0 {
		config.GatewayOpts = append([]gateway.ConfigOpt{
			gateway.WithURL(config.WsUrl),
		}, config.GatewayOpts...)

		config.Gateway = gateway.New(config.EventManager.HandleGatewayEvent, config.GatewayOpts...)
	}
	client.gateway = config.Gateway

	return client, nil
}

// FetchInstanceInfo fetches information about the instance.
func (c *clientImpl) FetchInstanceInfo() (models.InstanceInfo, error) {
	var res models.InstanceInfo
	_, err := c.request(InstanceInfo.Compile(nil), Data{}, &res)

	return res, err
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
