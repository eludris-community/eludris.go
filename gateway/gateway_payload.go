package gateway

import (
	"encoding/json"

	"github.com/eludris-community/eludris-api-types.go/v2/pandemonium"
)

type Payload struct {
	Op pandemonium.OpcodeType `json:"op"`
	D  any                    `json:"d,omitempty"`
}

type PayloadDataUnknown json.RawMessage

func (p *Payload) UnmarshalJSON(data []byte) error {
	var raw pandemonium.Payload
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var (
		payloadData any
		err         error
	)

	switch raw.Op {
	case pandemonium.PongOp:
	case pandemonium.RatelimitOp:
		var d pandemonium.Ratelimit
		err = json.Unmarshal(raw.D, &d)
		payloadData = d
	case pandemonium.HelloOp:
		var d pandemonium.Hello
		err = json.Unmarshal(raw.D, &d)
		payloadData = d
	case pandemonium.MessageCreateOp:
		var d pandemonium.MessageCreate
		err = json.Unmarshal(raw.D, &d)
		payloadData = d
	default:
		var d PayloadDataUnknown
		err = json.Unmarshal(raw.D, &d)
		payloadData = d
	}

	if err != nil {
		return err
	}
	p.Op = raw.Op
	p.D = payloadData
	return nil
}
