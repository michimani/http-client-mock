http-client-mock
===

[![codecov](https://codecov.io/gh/michimani/http-client-mock/branch/main/graph/badge.svg?token=MAS7YVKL9P)](https://codecov.io/gh/michimani/http-client-mock)

This is a package to generate a mock of `http.Client` in Go language. This package will help you to more easily write test code that mocks http request/response.

# Usage

```bash
go get github.com/michimani/http-client-mock/hcmock
```

## Sample

```go
import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/michimani/http-client-mock/hcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MyAPICall(t *testing.T) {
	cases := []struct {
		name   string
		method string
		in     *hcmock.MockInput
		expect http.Response
	}{
		{
			name:   "200 OK",
			method: http.MethodGet,
			in: &hcmock.MockInput{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Headers: []hcmock.Header{
					{Key: "header-1", Value: "value-1-1"},
					{Key: "header-1", Value: "value-1-2"},
					{Key: "header-2", Value: "value-2"},
				},
				BodyBytes: []byte("ok response"),
			},
			expect: http.Response{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"header-1": {"value-1-1", "value-1-2"},
					"header-2": {"value-2"},
				},
				Body:          io.NopCloser(bytes.NewReader([]byte("ok response"))),
				ContentLength: 11,
			},
		},
		{
			name:   "404 Not Found",
			method: http.MethodGet,
			in: &hcmock.MockInput{
				StatusCode: http.StatusNotFound,
				BodyBytes:  []byte("resource not found"),
			},
			expect: http.Response{
				Status:        "",
				StatusCode:    http.StatusNotFound,
				Header:        map[string][]string{},
				Body:          io.NopCloser(bytes.NewReader([]byte("resource not found"))),
				ContentLength: 18,
			},
		},
		{
			name:   "202 Accepted",
			method: http.MethodPost,
			in: &hcmock.MockInput{
				Status:     "202 Accepted",
				StatusCode: http.StatusAccepted,
			},
			expect: http.Response{
				Status:        "202 Accepted",
				StatusCode:    http.StatusAccepted,
				Header:        map[string][]string{},
				Body:          io.NopCloser(bytes.NewReader([]byte(""))),
				ContentLength: 0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			defer c.expect.Body.Close()

			asst := assert.New(tt)
			rqir := require.New(tt)

			mc := hcmock.New(c.in)

			req, err := http.NewRequest(c.method, "any-string-for-url", nil)
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
```

# License

[MIT](https://github.com/michimani/http-client-mock/blob/main/LICENSE)

# Author

[michimani210](https://twitter.com/michimani210)