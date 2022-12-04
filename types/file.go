// SPDX-License-Identifier: MIT

package types

// FileData represents the metadata of a file.
type FileData struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Bucket   string       `json:"bucket"`
	Spoiler  bool         `json:"spoiler"`
	Metadata FileMetadata `json:"metadata"`
}

// FileMetadata represents the metadata of a file's contents.
type FileMetadata struct {
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
