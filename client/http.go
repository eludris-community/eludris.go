package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ooliver1/eludris.go/types"
)

func (c clientImpl) request(method, path string, data any, obj any) error {
	payload, err := json.Marshal(data)

	if err != nil {
		return err
	}

	uri, err := url.Parse(fmt.Sprintf("%s/messages/", c.httpUrl))

	if err != nil {
		return err
	}

	fmt.Printf("Sending %s request to %s with payload %s\n", method, uri.String(), string(payload))

	req := http.Request{
		Method: method,
		URL:    uri,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: io.NopCloser(bytes.NewBuffer(payload)),
	}

	for {
		res, err := c.httpClient.Do(&req)

		if err != nil {
			return err
		}

		defer res.Body.Close()
		switch res.StatusCode {
		case 200:
			json.NewDecoder(res.Body).Decode(&obj)
			return nil
		case 429:
			var ratelimit types.RateLimit
			json.NewDecoder(res.Body).Decode(&ratelimit)
			retry_after := ratelimit.Data.RetryAfter
			time.Sleep(time.Duration(retry_after) * time.Millisecond)
		default:
			return fmt.Errorf("error sending message: %s", res.Status)
		}
	}
}
