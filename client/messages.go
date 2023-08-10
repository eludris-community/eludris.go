// SPDX-License-Identifier: MIT

package client

import "github.com/eludris-community/eludris-api-types.go/v2/models"

func (c clientImpl) SendMessage(author, content string) (models.Message, error) {
	msg := models.Message{Author: author, Content: content}

	var res models.Message
	_, err := c.request(SendMessage.Compile(nil), Data{Json: msg}, &res)

	return res, err
}
