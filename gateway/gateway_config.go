package gateway

import (
	"time"

	"github.com/gorilla/websocket"
)

func DefaultConfig() *Config {
	return &Config{}
}

type Config struct {
	URL          string
	Dialer       *websocket.Dialer
	PingInterval time.Duration
}

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

type ConfigOpt func(config *Config)

func WithURL(url string) ConfigOpt {
	return func(config *Config) {
		config.URL = url
	}
}

func WithDialer(dialer *websocket.Dialer) ConfigOpt {
	return func(config *Config) {
		config.Dialer = dialer
	}
}

func WithPingInterval(interval time.Duration) ConfigOpt {
	return func(config *Config) {
		config.PingInterval = interval
	}
}
