// SPDX-License-Identifier: MIT

package types

// Message represents a message sent to Eludris.
type Message struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}
