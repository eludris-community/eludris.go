// SPDX-License-Identifier: MIT

package events

import "github.com/eludris-community/eludris.go/v2/client"

func NewGenericEvent(client client.Client) *GenericEvent {
	return &GenericEvent{client: client}
}

type GenericEvent struct {
	client client.Client
}

func (e *GenericEvent) Client() client.Client {
	return e.client
}
