// SPDX-License-Identifier: MIT

package client

import "github.com/eludris-community/eludris-api-types.go/oprish"

// SendMessage sends a message to Eludris.
func (c clientImpl) SendMessage(author, content string) (oprish.Message, error) {
	msg := oprish.Message{Author: author, Content: content}

	var res oprish.Message
	_, err := c.request(Oprish, "POST", "/messages/", Data{Json: msg}, &res)

	return res, err
}
