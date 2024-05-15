package events

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
)

type UserCreate struct {
	*GenericEvent
	pandemonium.UserCreate
}

type UserUpdate struct {
	*GenericEvent
	pandemonium.UserUpdate
}

type PresenceUpdate struct {
	*GenericEvent
	pandemonium.PresenceUpdate
}
