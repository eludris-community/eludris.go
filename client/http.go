// SPDX-License-Identifier: MIT

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type RequestType int

const (
	Oprish RequestType = iota
	Effis
)

type Data struct {
	Json     any
	FormData map[string]io.Reader
}

func (c clientImpl) request(reqType RequestType, method, path string, data Data, obj any) (*http.Response, error) {
	// payload, err := json.Marshal(data)

	// if err != nil {
	// 	return err
	// }

	var base string
	switch reqType {
	case Oprish:
		base = c.httpUrl
	case Effis:
		base = c.fileUrl
	}

	uri, err := url.Parse(base + path)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Sending %s request to %s\n", method, uri.String())

	req := http.Request{
		Method: method,
		URL:    uri,
		Header: make(map[string][]string),
	}

	if data.Json != nil {
		req.Header.Set("Content-Type", "application/json")
		payload, err := json.Marshal(data.Json)

		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewBuffer(payload))
	} else if data.FormData != nil {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		for key, r := range data.FormData {
			var writer io.Writer
			if x, ok := r.(io.Closer); ok {
				defer x.Close()
			}
			// Add an image file
			if file, ok := r.(*os.File); ok {
				if writer, err = w.CreateFormFile(key, file.Name()); err != nil {
					return nil, err
				}
			} else {
				// Add other fields
				if writer, err = w.CreateFormField(key); err != nil {
					return nil, err
				}
			}
			if _, err = io.Copy(writer, r); err != nil {
				return nil, err
			}
		}

		w.Close()

		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Body = io.NopCloser(&b)
	}

	for {
		res, err := c.httpClient.Do(&req)

		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		switch res.StatusCode {
		case 200:
			json.NewDecoder(res.Body).Decode(&obj)
			return res, nil
		case 429:
			// TODO: Better ratelimiting using proper headers.
			retry_after, _ := strconv.Atoi(res.Header.Get("X-RateLimit-Reset"))
			time.Sleep(time.Duration(retry_after) * time.Millisecond)
		default:
			return nil, fmt.Errorf("error sending message: %s", res.Status)
		}
	}
}
