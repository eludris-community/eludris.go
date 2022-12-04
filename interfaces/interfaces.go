// SPDX-License-Identifier: MIT

// Package interfaces provides interfaces for the Eludris API.
package interfaces

import (
	"io"

	"github.com/eludris-community/eludris.go/types"
)

// Client represents a client for Eludris, with functions to interact with the API.
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
