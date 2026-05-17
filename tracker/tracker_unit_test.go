package tracker

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPushWebhook_Disabled verifies that the webhook is a no-op when disabled.
func TestPushWebhook_Disabled(t *testing.T) {
	tr := Init(Options{}, logger.Init(logger.Config{}))
	require.NoError(t, tr.PushWebhook(context.Background(), []byte("x"), nil))
}

// TestPushWebhook_Success exercises the happy path through an httptest server
// and asserts the body / headers are forwarded.
func TestPushWebhook_Success(t *testing.T) {
	var gotHeader string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Test")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	tr := Init(Options{Webhook: WebhookOptions{Enabled: true, URL: srv.URL}}, logger.Init(logger.Config{}))
	err := tr.PushWebhook(context.Background(), []byte("payload"), map[string]string{"X-Test": "yes"})
	require.NoError(t, err)
	assert.Equal(t, "yes", gotHeader)
	assert.Equal(t, []byte("payload"), gotBody)
}

// TestPushWebhook_Non200ReturnsError covers the non-OK branch.
func TestPushWebhook_Non200ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	tr := Init(Options{Webhook: WebhookOptions{Enabled: true, URL: srv.URL}}, logger.Init(logger.Config{}))
	err := tr.PushWebhook(context.Background(), []byte("payload"), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}
