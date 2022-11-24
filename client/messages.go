package client

import "github.com/ooliver1/eludris.go/types"

func (c clientImpl) SendMessage(author, content string) (types.Message, error) {
	msg := types.Message{Author: author, Content: content}

	var res types.Message
	err := c.request("POST", "/messages/", msg, &res)

	return res, err
}
