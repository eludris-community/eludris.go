// SPDX-License-Identifier: MIT

package types

type FileData struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Bucket   string       `json:"bucket"`
	Spoiler  bool         `json:"spoiler"`
	Metadata FileMetadata `json:"metadata"`
}

type FileMetadata struct {
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
