package rest

import "context"

func DefaultRequestConfig() *RequestConfig {
	return &RequestConfig{
		Ctx: context.TODO(),
	}
}

type RequestConfig struct {
	Ctx context.Context
}

type RequestOpt func(config *RequestConfig)

func (c *RequestConfig) Apply(opts []RequestOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithCtx(ctx context.Context) RequestOpt {
	return func(config *RequestConfig) {
		config.Ctx = ctx
	}
}
