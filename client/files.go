// SPDX-License-Identifier: MIT

package client

import (
	"io"
	"strconv"
	"strings"

	"github.com/eludris-community/eludris-api-types.go/v2/models"
)

func (c clientImpl) UploadAttachment(file io.Reader, spoiler bool) (models.FileData, error) {
	return c.UploadFile("attachments", file, spoiler)
}

func (c clientImpl) UploadFile(bucket string, file io.Reader, spoiler bool) (models.FileData, error) {
	var res models.FileData
	_, err := c.request(
		UploadFile.Compile(nil, bucket),
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
	res, err := c.request(FetchFile.Compile(nil, bucket, id), Data{}, nil)

	return res.Body, err
}

func (c clientImpl) FetchAttachmentData(id string) (models.FileData, error) {
	return c.FetchFileData("attachments", id)
}

func (c clientImpl) FetchFileData(bucket, id string) (models.FileData, error) {
	var res models.FileData
	_, err := c.request(FetchFileData.Compile(nil, bucket, id), Data{}, &res)

	return res, err
}

func (c clientImpl) FetchStaticFile(name string) (io.ReadCloser, error) {
	res, err := c.request(FetchStaticFile.Compile(nil, name), Data{}, nil)

	return res.Body, err
}
