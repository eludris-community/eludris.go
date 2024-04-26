// SPDX-License-Identifier: MIT

package rest

import (
	"github.com/eludris-community/eludris-api-types.go/v2/models"
	"github.com/eludris-community/eludris.go/v2/types"
)

func (r *restImpl) CreateMessage(messageCreate types.MessageCreate) (message models.Message, err error) {
	_, err = r.Request(SendMessage.Compile(nil), messageCreate, &message)
	return
}
