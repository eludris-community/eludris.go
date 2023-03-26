package client

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryValues map[string]any

type RequestType int

const (
	Oprish RequestType = iota
	Effis
)

func (q QueryValues) Encode() string {
	values := url.Values{}
	for k, v := range q {
		values.Set(k, fmt.Sprint(v))
	}
	return values.Encode()
}

type Endpoint struct {
	Method string
	Route  string
	Type   RequestType
}

type CompiledEndpoint struct {
	Endpoint *Endpoint
	Url      string
	Path     string
	Type     RequestType
}

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

	query := values.Encode()
	if query != "" {
		query = "?" + query
	}

	return &CompiledEndpoint{
		Endpoint: e,
		Url:      path + query,
		Path:     path,
		Type:     e.Type,
	}
}

func NewEndpoint(kind RequestType, method string, path string) *Endpoint {
	return &Endpoint{
		Method: method,
		Route:  path,
		Type:   kind,
	}
}
