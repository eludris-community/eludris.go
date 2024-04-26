// SPDX-License-Identifier: MIT

package rest

import "net/http"

// Files
var (
	// We do not need the attachments endpoint as that is just upload "attachment"
	// Upload
	UploadFile = NewEndpoint(Effis, http.MethodPost, "/<bucket>/")
	// Fetch
	GetFile = NewEndpoint(Effis, http.MethodGet, "/<bucket>/<id>")
	// Fetch Data
	GetFileData = NewEndpoint(Effis, http.MethodGet, "/<bucket>/<id>/data")
	// Download is not needed as that is just Fetch but `Content-Disposition: attachment`.
	// Static Files
	GetStaticFile = NewEndpoint(Effis, http.MethodGet, "/static/<name>")
)

// Instance
var (
	// Get Instance Info
	InstanceInfo = NewEndpoint(Oprish, http.MethodGet, "/")
)

// Messaging
var (
	// Create Message
	SendMessage = NewEndpoint(Oprish, http.MethodPost, "/messages/")
)
