package gqlclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Plug coverage gaps not exercised by client_test.go: cancelled context,
// non-multipart with files, request construction error, transport error,
// and the multipart form-data path.

func TestRun_CancelledContext(t *testing.T) {
	c := NewClient("http://example.com")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := c.Run(ctx, NewRequest("query{}"), nil)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestRun_FilesWithoutMultipartReturnsError(t *testing.T) {
	c := NewClient("http://example.com")
	req := NewRequest("query{}")
	req.File("upload", "x.txt", strings.NewReader("hi"))
	err := c.Run(context.Background(), req, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot send files")
}

func TestRun_HTTPNewRequestError(t *testing.T) {
	// Control characters in URL trigger http.NewRequest validation error.
	c := NewClient("http://exa\x00mple.com")
	err := c.Run(context.Background(), NewRequest("query{}"), nil)
	assert.Error(t, err)
}

// errRoundTripper always fails the request — exercises the Do error path
// for both runWithJSON and runWithPostFields.
type errRoundTripper struct{ err error }

func (e errRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, e.err
}

func TestRun_TransportError_JSON(t *testing.T) {
	c := NewClient("http://example.com", WithHTTPClient(&http.Client{Transport: errRoundTripper{err: errors.New("boom")}}))
	err := c.Run(context.Background(), NewRequest("query{}"), nil)
	assert.Error(t, err)
}

func TestRun_TransportError_Multipart(t *testing.T) {
	c := NewClient("http://example.com",
		UseMultipartForm(),
		WithHTTPClient(&http.Client{Transport: errRoundTripper{err: errors.New("boom")}}),
	)
	req := NewRequest("query{}")
	req.Var("k", "v")
	req.File("upload", "x.txt", strings.NewReader("hi"))
	err := c.Run(context.Background(), req, nil)
	assert.Error(t, err)
}

func TestRun_Multipart_Success(t *testing.T) {
	var hitMultipart bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			hitMultipart = true
		}
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, UseMultipartForm(), ImmediatelyCloseReqBody())
	req := NewRequest("query{}")
	req.Var("name", "jack")
	req.Header.Set("X-Custom", "hello")
	req.File("upload", "x.txt", strings.NewReader("file body"))
	var resp map[string]interface{}
	err := c.Run(context.Background(), req, &resp)
	assert.NoError(t, err)
	assert.True(t, hitMultipart)
	assert.Equal(t, true, resp["ok"])
}

func TestRun_Multipart_NoFiles_NoVars(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, UseMultipartForm())
	err := c.Run(context.Background(), NewRequest("query{}"), nil)
	assert.NoError(t, err)
}

func TestRun_JSON_Non200WithBadBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL)
	err := c.Run(context.Background(), NewRequest("query{}"), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-200")
}

func TestRequest_Accessors(t *testing.T) {
	req := NewRequest("query{}")
	req.Var("a", 1)
	req.File("file", "x.txt", strings.NewReader("body"))
	assert.Equal(t, "query{}", req.Query())
	assert.Equal(t, map[string]interface{}{"a": 1}, req.Vars())
	assert.Len(t, req.Files(), 1)
}

func TestClientOptions_ImmediatelyCloseReqBody(t *testing.T) {
	c := NewClient("http://x", ImmediatelyCloseReqBody()).(*Client)
	assert.True(t, c.closeReq)
}

func TestClientOptions_DefaultHTTPClientIsSet(t *testing.T) {
	c := NewClient("http://x").(*Client)
	assert.NotNil(t, c.httpClient)
}

func TestGraphErr_Error(t *testing.T) {
	err := graphErr{Message: "bad"}
	assert.Equal(t, "graphql: bad", err.Error())
}

// failingReader returns an error on the first Read so that io.Copy fails.
type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("read fail") }

// Triggers the "preparing file" io.Copy error branch in runWithPostFields.
func TestRun_Multipart_FilePrepareError(t *testing.T) {
	c := NewClient("http://example.com",
		UseMultipartForm(),
		WithHTTPClient(&http.Client{Transport: errRoundTripper{err: errors.New("never reached")}}),
	)
	req := NewRequest("query{}")
	req.File("upload", "x.txt", failingReader{})
	err := c.Run(context.Background(), req, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "preparing file")
}

// Triggers the "encode variables" json error branch in runWithPostFields by
// passing a value json cannot marshal (channels are not encodable).
func TestRun_Multipart_EncodeVariablesError(t *testing.T) {
	c := NewClient("http://example.com",
		UseMultipartForm(),
		WithHTTPClient(&http.Client{Transport: errRoundTripper{err: errors.New("never reached")}}),
	)
	req := NewRequest("query{}")
	req.Var("ch", make(chan int))
	err := c.Run(context.Background(), req, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encode variables")
}

// Triggers the "encode body" json error branch in runWithJSON.
func TestRun_JSON_EncodeBodyError(t *testing.T) {
	c := NewClient("http://example.com",
		WithHTTPClient(&http.Client{Transport: errRoundTripper{err: errors.New("never reached")}}),
	)
	req := NewRequest("query{}")
	req.Var("ch", make(chan int))
	err := c.Run(context.Background(), req, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encode body")
}

func TestRedactHeaders(t *testing.T) {
	in := http.Header{
		"Authorization": []string{"Bearer secret-token"},
		"Cookie":        []string{"session=abc"},
		"X-Api-Key":     []string{"key-value"},
		"User-Agent":    []string{"unit-test"},
	}
	got := redactHeaders(in)

	assert.Equal(t, []string{"[REDACTED]"}, got["Authorization"])
	assert.Equal(t, []string{"[REDACTED]"}, got["Cookie"])
	assert.Equal(t, []string{"[REDACTED]"}, got["X-Api-Key"])
	assert.Equal(t, []string{"unit-test"}, got["User-Agent"])
	// Original must not have been mutated.
	assert.Equal(t, "Bearer secret-token", in.Get("Authorization"))
}

func TestTruncateForLog(t *testing.T) {
	short := strings.Repeat("a", maxLoggedBodyBytes)
	long := strings.Repeat("b", maxLoggedBodyBytes+50)

	assert.Equal(t, short, truncateForLog(short))
	got := truncateForLog(long)
	assert.True(t, strings.HasSuffix(got, "(truncated)"))
	assert.Equal(t, maxLoggedBodyBytes+len("(truncated)"), len(got))
}

// Triggers the multipart non-200 + bad JSON body path.
func TestRun_Multipart_Non200WithBadBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, UseMultipartForm())
	err := c.Run(context.Background(), NewRequest("query{}"), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-200")
}
