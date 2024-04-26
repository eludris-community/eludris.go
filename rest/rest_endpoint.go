// SPDX-License-Identifier: MIT

package rest

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryValues map[string]any

// RequestType represents the destination for the request - oprish or effis.
type RequestType int

const (
	Oprish RequestType = iota
	Effis
)

// Encode the query values into a string via `url.Values`.
func (q QueryValues) Encode() string {
	values := url.Values{}
	for k, v := range q {
		values.Set(k, fmt.Sprint(v))
	}
	return values.Encode()
}

// Endpoint represents a single endpoint.
// It contains all the information needed to construct a request from input.
type Endpoint struct {
	Method string
	Route  string
	Type   RequestType
}

// CompiledEndpoint represents a compiled endpoint.
// It contains the endpoint, the constructed URL, the path, and the type.
type CompiledEndpoint struct {
	Endpoint *Endpoint
	Url      string
	Path     string
	Type     RequestType
}

// Compile the endpoint with the given values and parameters.
func (e *Endpoint) Compile(values QueryValues, params ...any) *CompiledEndpoint {
	path := e.Route
	for _, param := range params {
		// Replace the first occurence of <...> with the param.
		start := strings.Index(path, "<")
		end := strings.Index(path, ">")
		if start == -1 || end == -1 {
			break
		}
		paramValue := fmt.Sprint(param)
		path = path[:start] + paramValue + path[end+1:]
	}

	var query string
	if values == nil {
		query = ""
	} else {
		query = values.Encode()
		if query != "" {
			query = "?" + query
		}
	}

	return &CompiledEndpoint{
		Endpoint: e,
		Url:      path + query,
		Path:     path,
		Type:     e.Type,
	}
}

// NewEndpoint creates a new endpoint from the given values.
func NewEndpoint(kind RequestType, method string, path string) *Endpoint {
	return &Endpoint{
		Method: method,
		Route:  path,
		Type:   kind,
	}
}
