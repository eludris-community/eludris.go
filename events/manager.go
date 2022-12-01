package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/ooliver1/eludris.go/interfaces"
)

type EventListener interface {
	Handle(Event)
	Op() string
	Func() func(Event)
}

type eventListener[E Event] struct {
	op string
	f  func(E)
}

func (e eventListener[E]) Op() string {
	return e.op
}

func (e eventListener[E]) Func() func(Event) {
	// THe pain behind it all.
	return any(e.f).(func(Event))
}

func (l eventListener[E]) Handle(event Event) {
	if event, ok := event.(E); ok {
		l.f(event)
	}
}

type EventManager interface {
	Subscribe(EventListener)
	Dispatch(interfaces.Client, []byte)
}

type managerImpl struct {
	subscribers map[string][]EventListener
}

func Subscribe[E Event](m EventManager, subscriber func(E)) {
	t := reflect.TypeOf(subscriber)
	eventT := t.In(0)

	for i := 0; i < eventT.NumMethod(); i++ {
		method := eventT.Method(i)
		if method.Name == "Op" {
			op := method.Func.Call([]reflect.Value{reflect.New(eventT).Elem()})[0].String()
			listener := eventListener[E]{op: op, f: subscriber}
			m.Subscribe(&listener)
			return
		}
	}
}

func (m *managerImpl) Subscribe(listener EventListener) {
	if _, ok := m.subscribers[listener.Op()]; !ok {
		m.subscribers[listener.Op()] = []EventListener{}
	}

	m.subscribers[listener.Op()] = append(m.subscribers[listener.Op()], listener)
}

func (m *managerImpl) Dispatch(client interfaces.Client, data []byte) {
	var msg map[string]any
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&msg)

	op := msg["op"].(string)
	innerData := msg["d"]
	if innerData == nil {
        innerData = make(map[string]any)
    }

	subscribers := m.subscribers[op]

	var event Event

	switch op {

	case "MESSAGE_CREATE":
		{
			var innerEvent MessageEvent
			mapstructure.Decode(innerData, &innerEvent)
			event = &innerEvent
		}
	case "PONG":
		{
			var innerEvent PongEvent
			mapstructure.Decode(innerData, &innerEvent)
			event = &innerEvent
		}
	default:
		fmt.Println("Unknown event type: ", op)

		return
	}

	event.SetClient(client)

	for _, subscriber := range subscribers {
		subscriber.Handle(event)
	}
}

func NewEventManager() EventManager {
	return &managerImpl{subscribers: make(map[string][]EventListener)}
}
