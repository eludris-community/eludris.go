package eludris

import (
	"github.com/eludris-community/eludris.go/v2/client"
	"github.com/eludris-community/eludris.go/v2/handlers"
)

// New creates a new client.
func New(token string, opts ...client.ConfigOpt) (client.Client, error) {
	config := client.DefaultConfig(handlers.AllHandlers())
	config.Apply(opts)

	return client.BuildClient(token, config)
}
