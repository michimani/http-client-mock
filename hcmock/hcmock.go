package hcmock

import (
	"bytes"
	"io"
	"net/http"
)

var _ http.RoundTripper = (*roundTripFunc)(nil)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type Header struct {
	Key   string
	Value string
}

// MockInput is input for generating http.Client mock.
type MockInput struct {
	Status     string
	StatusCode int
	Headers    []Header
	BodyBytes  []byte
}

// New returns a pointer to `http.Client` that returns the status code, header, and body specified in MockInput
// as the `http.Response` when the `Do()` method is executed.
// If MockInput is nil, New returns `http.DefaultClient`.
func New(in *MockInput) *http.Client {
	if in == nil {
		return http.DefaultClient
	}

	header := map[string][]string{}
	for _, h := range in.Headers {
		if len(header[h.Key]) > 0 {
			header[h.Key] = append(header[h.Key], h.Value)
		} else {
			header[h.Key] = []string{h.Value}
		}
	}

	return newMockClient(func(req *http.Request) *http.Response {
		return &http.Response{
			Status:        in.Status,
			StatusCode:    in.StatusCode,
			Body:          io.NopCloser(bytes.NewReader(in.BodyBytes)),
			Header:        header,
			ContentLength: int64(len(in.BodyBytes)),
		}
	})
}

func newMockClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}
