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

	"github.com/apex/log"
)

type Data struct {
	Json     any
	FormData map[string]io.Reader
}

func (c clientImpl) retry(endpoint *CompiledEndpoint, data Data, tries int, obj any) (*http.Response, error) {
	var base string
	switch endpoint.Type {
	case Oprish:
		base = c.httpUrl
	case Effis:
		base = c.fileUrl
	}

	uri, err := url.Parse(base + endpoint.Path)

	if err != nil {
		return nil, err
	}

	method := endpoint.Endpoint.Method

	req := &http.Request{
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

	first, err := c.rateLimiter.WaitBucket(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error waiting for bucket: %w", err)
	}

	log.WithField("method", method).WithField("url", uri.String()).Debug("Sending request")
	res, err := c.httpClient.Do(req)

	if err != nil {
		_ = c.rateLimiter.UnlockBucket(endpoint, nil, first)
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if err := c.rateLimiter.UnlockBucket(endpoint, res, first); err != nil {
		return nil, fmt.Errorf("error unlocking bucket: %w", err)
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		json.NewDecoder(res.Body).Decode(&obj)
		return res, nil
	case 429:
		log.WithField("path", endpoint.Path).Warn("Rate limit exceeded")
		if tries >= c.rateLimiter.MaxRetries() {
			return nil, fmt.Errorf("rate limit exceeded")
		}

		return c.retry(endpoint, data, tries+1, obj)
	default:
		return nil, fmt.Errorf("error sending message: %s", res.Status)
	}
}

func (c clientImpl) request(endpoint *CompiledEndpoint, data Data, obj any) (*http.Response, error) {
	return c.retry(endpoint, data, 1, obj)
}
