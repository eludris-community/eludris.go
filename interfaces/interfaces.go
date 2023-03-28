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

	// Send a message to the chat.
	// Returns the message that was sent.
	// Author must be between 1 and 32 characters long.
	// Content must be between 1 and the instance's configured max message length.
	SendMessage(author, content string) (models.Message, error)

	// Upload a file to the instance's attachments bucket.
	UploadAttachment(file io.Reader, spoiler bool) (models.FileData, error)
	// Upload a file to a specific bucket, currently only "attachments" exists.
	UploadFile(bucket string, file io.Reader, spoiler bool) (models.FileData, error)
	// Fetch a file from the instance's attachments bucket.
	FetchAttachment(id string) (io.ReadCloser, error)
	// Fetch a file from a specific bucket, currently only "attachments" exists.
	FetchFile(bucket, id string) (io.ReadCloser, error)
	// Fetch a file from the instance's attachments bucket.
	FetchAttachmentData(id string) (models.FileData, error)
	// Fetch a file from a specific bucket, currently only "attachments" exists.
	FetchFileData(bucket, id string) (models.FileData, error)
}
