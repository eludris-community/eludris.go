// SPDX-License-Identifier: MIT

package interfaces

import (
	"io"

	"github.com/eludris-community/eludris.go/types"
)

type Client interface {
	Connect() error

	SendMessage(author, content string) (types.Message, error)

	UploadAttachment(file io.Reader, spoiler bool) (types.FileData, error)
	UploadFile(bucket string, file io.Reader, spoiler bool) (types.FileData, error)
	FetchAttachment(id string) (io.ReadCloser, error)
	FetchFile(bucket, id string) (io.ReadCloser, error)
	FetchAttachmentData(id string) (types.FileData, error)
	FetchFileData(bucket, id string) (types.FileData, error)
}
