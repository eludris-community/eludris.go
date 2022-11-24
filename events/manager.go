package events

import (
	"bytes"
	"encoding/json"

	"github.com/ooliver1/eludris.go/interfaces"
)

type EventManager interface {
	Subscribe(func(MessageEvent))
	Dispatch(interfaces.Client, []byte)
}

type managerImpl struct {
	subscribers []func(MessageEvent)
}

func (m *managerImpl) Subscribe(subscriber func(MessageEvent)) {
	// t := reflect.TypeOf(subscriber)
	// actualFuckingType := t.In(0)
	m.subscribers = append(m.subscribers, subscriber)
}

func (m *managerImpl) Dispatch(client interfaces.Client, data []byte) {
	var msg MessageEvent
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&msg)

	msg.SetClient(client)

	for _, subscriber := range m.subscribers {
		subscriber(msg)
	}
}

func NewEventManager() EventManager {
	return &managerImpl{}
}
