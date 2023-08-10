package events

import (
	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
)

type MessageCreate struct {
	*GenericEvent
	pandemonium.MessageCreate
}
