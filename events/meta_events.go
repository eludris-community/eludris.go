package events

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
)

type Pong struct {
	*GenericEvent
}

type Hello struct {
	*GenericEvent
	pandemonium.Hello
}

type RateLimit struct {
	*GenericEvent
	pandemonium.Ratelimit
}

type Authenticated struct {
	*GenericEvent
	pandemonium.Authenticated
}
