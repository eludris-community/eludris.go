// SPDX-License-Identifier: MIT

package types

type RateLimit struct {
	Data RateLimitData `json:"data"`
}

type RateLimitData struct {
	RetryAfter uint `json:"retry_after"`
}
