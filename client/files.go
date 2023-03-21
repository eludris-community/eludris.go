// SPDX-License-Identifier: MIT

package client

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/eludris-community/eludris-api-types.go/effis"
)

// UploadAttachment uploads an attachment to the file server.
func (c clientImpl) UploadAttachment(file io.Reader, spoiler bool) (effis.FileData, error) {
	return c.UploadFile("attachments", file, spoiler)
}

// UploadStaticFile uploads a file to the file server with the chosen "bucket".
func (c clientImpl) UploadFile(bucket string, file io.Reader, spoiler bool) (effis.FileData, error) {
	var res effis.FileData
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

// FetchAttachment fetches the raw data of an attachment.
func (c clientImpl) FetchAttachment(id string) (io.ReadCloser, error) {
	return c.FetchFile("attachments", id)
}

// FetchFile fetches the raw data of a file.
func (c clientImpl) FetchFile(bucket, id string) (io.ReadCloser, error) {
	res, err := c.request(Effis, "GET", fmt.Sprintf("/%s/%s", bucket, id), Data{}, nil)

	return res.Body, err
}

// FetchAttachmentData fetches the metadata of an attachment.
func (c clientImpl) FetchAttachmentData(id string) (effis.FileData, error) {
	return c.FetchFileData("attachments", id)
}

// FetchFileData fetches the metadata of a file.
func (c clientImpl) FetchFileData(bucket, id string) (effis.FileData, error) {
	var res effis.FileData
	_, err := c.request(Effis, "GET", fmt.Sprintf("/%s/%s", bucket, id), Data{}, &res)

	return res, err
}

// FetchStaticFile fetches the raw data of a static file.
func (c clientImpl) FetchStaticFile(name string) (io.ReadCloser, error) {
	res, err := c.request(Effis, "GET", fmt.Sprintf("/static/%s", name), Data{}, nil)

	return res.Body, err
}
