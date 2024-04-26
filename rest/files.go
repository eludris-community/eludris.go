// SPDX-License-Identifier: MIT

package rest

import (
	"io"
	"strconv"
	"strings"

	"github.com/eludris-community/eludris-api-types.go/v2/models"
)

func (r *restImpl) UploadAttachment(file io.Reader, spoiler bool) (models.FileData, error) {
	return r.UploadFile("attachments", file, spoiler)
}

func (r *restImpl) UploadFile(bucket string, fileData io.Reader, spoiler bool) (file models.FileData, err error) {
	_, err = r.Request(
		UploadFile.Compile(nil, bucket),
		Form(map[string]io.Reader{
			"file":    fileData,
			"spoiler": strings.NewReader(strconv.FormatBool(spoiler)),
		}),
		&file,
	)
	return
}

func (r *restImpl) GetAttachment(id string) (io.ReadCloser, error) {
	return r.GetFile("attachments", id)
}

func (r *restImpl) GetFile(bucket, id string) (io.ReadCloser, error) {
	res, err := r.Request(GetFile.Compile(nil, bucket, id), nil, nil)
	return res.Body, err
}

func (r *restImpl) GetAttachmentData(id string) (models.FileData, error) {
	return r.GetFileData("attachments", id)
}

func (r *restImpl) GetFileData(bucket, id string) (fileData models.FileData, err error) {
	_, err = r.Request(GetFileData.Compile(nil, bucket, id), nil, &fileData)
	return
}

func (r *restImpl) GetStaticFile(name string) (io.ReadCloser, error) {
	res, err := r.Request(GetStaticFile.Compile(nil, name), nil, nil)
	return res.Body, err
}
