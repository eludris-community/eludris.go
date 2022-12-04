// SPDX-License-Identifier: MIT

package client

import "github.com/eludris-community/eludris.go/types"

// SendMessage sends a message to Eludris.
func (c clientImpl) SendMessage(author, content string) (types.Message, error) {
	msg := types.Message{Author: author, Content: content}

	var res types.Message
	_, err := c.request(Oprish, "POST", "/messages/", Data{Json: msg}, &res)

	return res, err
}
