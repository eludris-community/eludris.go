// SPDX-License-Identifier: MIT

package client

import (
	"io"
	"strconv"
	"strings"

	"github.com/eludris-community/eludris-api-types.go/models"
)

// UploadAttachment uploads an attachment to the file server.
func (c clientImpl) UploadAttachment(file io.Reader, spoiler bool) (models.FileData, error) {
	return c.UploadFile("attachments", file, spoiler)
}

// UploadStaticFile uploads a file to the file server with the chosen "bucket".
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

// FetchAttachment fetches the raw data of an attachment.
func (c clientImpl) FetchAttachment(id string) (io.ReadCloser, error) {
	return c.FetchFile("attachments", id)
}

// FetchFile fetches the raw data of a file.
func (c clientImpl) FetchFile(bucket, id string) (io.ReadCloser, error) {
	res, err := c.request(FetchFile.Compile(nil, bucket, id), Data{}, nil)

	return res.Body, err
}

// FetchAttachmentData fetches the metadata of an attachment.
func (c clientImpl) FetchAttachmentData(id string) (models.FileData, error) {
	return c.FetchFileData("attachments", id)
}

// FetchFileData fetches the metadata of a file.
func (c clientImpl) FetchFileData(bucket, id string) (models.FileData, error) {
	var res models.FileData
	_, err := c.request(FetchFileData.Compile(nil, bucket, id), Data{}, &res)

	return res, err
}

// FetchStaticFile fetches the raw data of a static file.
func (c clientImpl) FetchStaticFile(name string) (io.ReadCloser, error) {
	res, err := c.request(FetchStaticFile.Compile(nil, name), Data{}, nil)

	return res.Body, err
}
