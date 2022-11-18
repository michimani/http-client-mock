package hcmock_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michimani/http-client-mock/hcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MockClientDo(t *testing.T) {
	cases := []struct {
		name   string
		in     *hcmock.MockInput
		expect http.Response
	}{
		{
			name: "ok: all",
			in: &hcmock.MockInput{
				Status:     "test status",
				StatusCode: http.StatusOK,
				Headers: []hcmock.Header{
					{Key: "header-1", Value: "value-1-1"},
					{Key: "header-1", Value: "value-1-2"},
					{Key: "header-2", Value: "value-2"},
				},
				BodyBytes: []byte("test response body"),
			},
			expect: http.Response{
				Status:     "test status",
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"header-1": {"value-1-1", "value-1-2"},
					"header-2": {"value-2"},
				},
				Body:          io.NopCloser(bytes.NewReader([]byte("test response body"))),
				ContentLength: 18,
			},
		},
		{
			name: "ok: no status",
			in: &hcmock.MockInput{
				StatusCode: http.StatusOK,
				Headers: []hcmock.Header{
					{Key: "header-1", Value: "value-1-1"},
					{Key: "header-1", Value: "value-1-2"},
					{Key: "header-2", Value: "value-2"},
				},
				BodyBytes: []byte("test response body"),
			},
			expect: http.Response{
				Status:     "",
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"header-1": {"value-1-1", "value-1-2"},
					"header-2": {"value-2"},
				},
				Body:          io.NopCloser(bytes.NewReader([]byte("test response body"))),
				ContentLength: 18,
			},
		},
		{
			name: "ok: no header",
			in: &hcmock.MockInput{
				Status:     "test status",
				StatusCode: http.StatusOK,
				BodyBytes:  []byte("test response body"),
			},
			expect: http.Response{
				Status:        "test status",
				StatusCode:    http.StatusOK,
				Header:        map[string][]string{},
				Body:          io.NopCloser(bytes.NewReader([]byte("test response body"))),
				ContentLength: 18,
			},
		},
		{
			name: "ok: no response body",
			in: &hcmock.MockInput{
				Status:     "test status",
				StatusCode: http.StatusOK,
				Headers: []hcmock.Header{
					{Key: "header-1", Value: "value-1-1"},
					{Key: "header-1", Value: "value-1-2"},
					{Key: "header-2", Value: "value-2"},
				},
			},
			expect: http.Response{
				Status:     "test status",
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"header-1": {"value-1-1", "value-1-2"},
					"header-2": {"value-2"},
				},
				Body:          io.NopCloser(bytes.NewReader([]byte(""))),
				ContentLength: 0,
			},
		},
		{
			name: "ok: MockInput is nil (default client)",
			in:   nil,
			expect: http.Response{
				Status:        "200 OK",
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(bytes.NewReader([]byte("test"))),
				ContentLength: 4,
			},
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("test"))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			defer c.expect.Body.Close()

			asst := assert.New(tt)
			rqir := require.New(tt)

			mc := hcmock.New(c.in)

			req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			rqir.NoError(err)
			rqir.NotNil(req)

			res, err := mc.Do(req)
			rqir.NoError(err)
			rqir.NotNil(res)
			defer res.Body.Close()

			asst.Equal(c.expect.Status, res.Status)
			asst.Equal(c.expect.StatusCode, res.StatusCode)
			asst.Equal(c.expect.ContentLength, res.ContentLength)

			eb := new(bytes.Buffer)
			_, err = io.Copy(eb, c.expect.Body)
			rqir.NoError(err)

			rb := new(bytes.Buffer)
			_, err = io.Copy(rb, res.Body)
			rqir.NoError(err)

			asst.Equal(eb, rb)
		})
	}
}
