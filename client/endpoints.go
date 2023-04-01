// SPDX-License-Identifier: MIT

package client

// Oprish

// Instance Info
var (
	InstanceInfo = NewEndpoint(Oprish, "GET", "/")
)

// Create Message
var (
	SendMessage = NewEndpoint(Oprish, "POST", "/messages/")
)

// Effis
// We do not need the attachments endpoint as that is just upload "attachment"
// Upload
var (
	UploadFile = NewEndpoint(Effis, "POST", "/<bucket>/")
)

// Fetch
var (
	FetchFile = NewEndpoint(Effis, "GET", "/<bucket>/<id>")
)

// Fetch Data
var (
	FetchFileData = NewEndpoint(Effis, "GET", "/<bucket>/<id>/data")
)

// Download is not needed as that is just Fetch but `Content-Disposition: attachment`.

// Static Files
var (
	FetchStaticFile = NewEndpoint(Effis, "GET", "/static/<name>")
)
