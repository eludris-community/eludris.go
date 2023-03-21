// SPDX-License-Identifier: MIT

// Package interfaces provides interfaces for the Eludris API.
package interfaces

import (
	"io"

	"github.com/eludris-community/eludris-api-types.go/effis"
	"github.com/eludris-community/eludris-api-types.go/oprish"
)

// Client represents a client for Eludris, with functions to interact with the API.
type Client interface {
	Connect() error

	SendMessage(author, content string) (oprish.Message, error)

	UploadAttachment(file io.Reader, spoiler bool) (effis.FileData, error)
	UploadFile(bucket string, file io.Reader, spoiler bool) (effis.FileData, error)
	FetchAttachment(id string) (io.ReadCloser, error)
	FetchFile(bucket, id string) (io.ReadCloser, error)
	FetchAttachmentData(id string) (effis.FileData, error)
	FetchFileData(bucket, id string) (effis.FileData, error)
}
