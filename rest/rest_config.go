package rest

import "net/http"

type Config struct {
	HttpClient      *http.Client
	RateLimiter     RateLimiter
	RateLimiterOpts []RateLimiterOpt
	ApiUrl          string
	FileUrl         string
}

func DefaultConfig(apiUrl, fileUrl string) *Config {
	return &Config{
		HttpClient: &http.Client{},
		ApiUrl:     apiUrl,
		FileUrl:    fileUrl,
	}
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithHttpClient(httpClient *http.Client) ConfigOpt {
	return func(config *Config) {
		config.HttpClient = httpClient
	}
}

func WithRateLimiter(rateLimiter RateLimiter) ConfigOpt {
	return func(config *Config) {
		config.RateLimiter = rateLimiter
	}
}

func WithRateLimiterOpts(opts ...RateLimiterOpt) ConfigOpt {
	return func(config *Config) {
		config.RateLimiterOpts = append(config.RateLimiterOpts, opts...)
	}
}

func WithApiUrl(apiUrl string) ConfigOpt {
	return func(config *Config) {
		config.ApiUrl = apiUrl
	}
}

func WithFileUrl(fileUrl string) ConfigOpt {
	return func(config *Config) {
		config.FileUrl = fileUrl
	}
}
