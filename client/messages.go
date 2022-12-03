package client

import "github.com/ooliver1/eludris.go/types"

func (c clientImpl) SendMessage(author, content string) (types.Message, error) {
	msg := types.Message{Author: author, Content: content}

	var res types.Message
	_, err := c.request(Oprish, "POST", "/messages/", Data{Json: msg}, &res)

	return res, err
}
