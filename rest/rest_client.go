// SPDX-License-Identifier: MIT

package rest

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
	"github.com/eludris-community/eludris-api-types.go/v2/models"
	"github.com/eludris-community/eludris.go/v2/types"
)

func New(apiUrl, fileUrl string, opts ...ConfigOpt) Rest {
	config := DefaultConfig(apiUrl, fileUrl)
	config.Apply(opts)

	if config.RateLimiter == nil {
		config.RateLimiter = NewRateLimiter(config.RateLimiterOpts...)
	}

	return &restImpl{
		config: *config,
	}
}

type Rest interface {
	Request(endpoint *CompiledEndpoint, data any, response any, opts ...RequestOpt) (*http.Response, error)
	HttpClient() *http.Client
	RateLimiter() RateLimiter
	// Files
	UploadAttachment(file io.Reader, spoiler bool) (models.FileData, error)
	UploadFile(bucket string, fileData io.Reader, spoiler bool) (file models.FileData, err error)
	GetAttachment(id string) (io.ReadCloser, error)
	GetFile(bucket, id string) (io.ReadCloser, error)
	GetAttachmentData(id string) (models.FileData, error)
	GetFileData(bucket, id string) (fileData models.FileData, err error)
	GetStaticFile(name string) (io.ReadCloser, error)
	// Instance
	GetInstanceInfo() (info models.InstanceInfo, err error)
	// Messaging
	CreateMessage(messageCreate types.MessageCreate) (message models.Message, err error)
}

type restImpl struct {
	config Config
}

func (r *restImpl) HttpClient() *http.Client {
	return r.config.HttpClient
}

func (r *restImpl) RateLimiter() RateLimiter {
	return r.config.RateLimiter
}

type Form map[string]io.Reader

func (r *restImpl) retry(endpoint *CompiledEndpoint, data any, tries int, response any, opts ...RequestOpt) (*http.Response, error) {
	var base string
	switch endpoint.Type {
	case Oprish:
		base = r.config.ApiUrl
	case Effis:
		base = r.config.FileUrl
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

	if data != nil {
		switch data.(type) {
		case Form:
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
			for key, r := range data.(Form) {
				var writer io.Writer
				if x, ok := r.(io.Closer); ok {
					defer x.Close()
				}

				if file, ok := r.(*os.File); ok {
					if writer, err = w.CreateFormFile(key, file.Name()); err != nil {
						return nil, err
					}
				} else {
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
		default:
			req.Header.Set("Content-Type", "application/json")
			payload, err := json.Marshal(data)

			if err != nil {
				return nil, err
			}

			req.Body = io.NopCloser(bytes.NewBuffer(payload))
		}
	}

	config := DefaultRequestConfig()
	config.Apply(opts)

	first, err := r.RateLimiter().WaitBucket(config.Ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("error waiting for bucket: %w", err)
	}

	log.WithField("method", method).WithField("url", uri.String()).Debug("Sending request")
	res, err := r.HttpClient().Do(req)

	if err != nil {
		_ = r.RateLimiter().UnlockBucket(config.Ctx, endpoint, nil, first)
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if err := r.RateLimiter().UnlockBucket(config.Ctx, endpoint, res, first); err != nil {
		return nil, fmt.Errorf("error unlocking bucket: %w", err)
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		json.NewDecoder(res.Body).Decode(&response)
		return res, nil
	case http.StatusTooManyRequests:
		log.WithField("path", endpoint.Path).Warn("Rate limit exceeded")
		if tries >= r.RateLimiter().MaxRetries() {
			return nil, fmt.Errorf("rate limit exceeded")
		}

		return r.retry(endpoint, data, tries+1, response)
	default:
		return nil, fmt.Errorf("error sending message: %s", res.Status)
	}
}

func (r *restImpl) Request(endpoint *CompiledEndpoint, data any, response any, opts ...RequestOpt) (*http.Response, error) {
	return r.retry(endpoint, data, 1, response, opts...)
}
