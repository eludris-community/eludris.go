package interfaces

import (
	"github.com/ooliver1/eludris.go/types"
)

type Client interface {
	SendMessage(message, author string) (types.Message, error)
	Connect() error
}
