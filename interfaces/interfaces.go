package interfaces

import (
	"github.com/ooliver1/eludris.go/types"
)

type Client interface {
	SendMessage(author, content string) (types.Message, error)
	Connect() error
}
