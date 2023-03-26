// SPDX-License-Identifier: MIT

// Package interfaces provides interfaces for the Eludris API.
package interfaces

import (
	"io"

	"github.com/eludris-community/eludris-api-types.go/models"
)

// Client represents a client for Eludris, with functions to interact with the API.
type Client interface {
	Connect() error

	SendMessage(author, content string) (models.Message, error)

	UploadAttachment(file io.Reader, spoiler bool) (models.FileData, error)
	UploadFile(bucket string, file io.Reader, spoiler bool) (models.FileData, error)
	FetchAttachment(id string) (io.ReadCloser, error)
	FetchFile(bucket, id string) (io.ReadCloser, error)
	FetchAttachmentData(id string) (models.FileData, error)
	FetchFileData(bucket, id string) (models.FileData, error)
}
