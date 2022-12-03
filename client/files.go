// SPDX-License-Identifier: MIT

package client

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/eludris-community/eludris.go/types"
)

func (c clientImpl) UploadAttachment(file io.Reader, spoiler bool) (types.FileData, error) {
	return c.UploadFile("attachments", file, spoiler)
}

func (c clientImpl) UploadFile(bucket string, file io.Reader, spoiler bool) (types.FileData, error) {
	var res types.FileData
	_, err := c.request(
		Effis,
		"POST",
		fmt.Sprintf("/%s/", bucket),
		Data{FormData: map[string]io.Reader{
			"file":    file,
			"spoiler": strings.NewReader(strconv.FormatBool(spoiler)),
		}},
		&res,
	)

	return res, err
}

func (c clientImpl) FetchAttachment(id string) (io.ReadCloser, error) {
	return c.FetchFile("attachments", id)
}

func (c clientImpl) FetchFile(bucket, id string) (io.ReadCloser, error) {
	res, err := c.request(Effis, "GET", fmt.Sprintf("/%s/%s", bucket, id), Data{}, nil)

	return res.Body, err
}

func (c clientImpl) FetchAttachmentData(id string) (types.FileData, error) {
	return c.FetchFileData("attachments", id)
}

func (c clientImpl) FetchFileData(bucket, id string) (types.FileData, error) {
	var res types.FileData
	_, err := c.request(Effis, "GET", fmt.Sprintf("/%s/%s", bucket, id), Data{}, &res)

	return res, err
}

func (c clientImpl) FetchStaticFile(name string) (io.ReadCloser, error) {
	res, err := c.request(Effis, "GET", fmt.Sprintf("/static/%s", name), Data{}, nil)

	return res.Body, err
}
