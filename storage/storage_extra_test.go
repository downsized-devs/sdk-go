package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeS3 implements the s3API seam. Each operation records its input and
// returns the canned output/error configured on the struct.
type fakeS3 struct {
	putInput *s3.PutObjectInput
	putErr   error

	getInput *s3.GetObjectInput
	getBody  []byte
	getErr   error
	// getBodyOverride, if non-nil, is returned as the GetObjectOutput body
	// instead of wrapping getBody in a NopCloser.
	getBodyOverride io.ReadCloser
	// getContentLengthOverride, if non-nil, replaces the auto-derived
	// ContentLength on the GetObjectOutput.
	getContentLengthOverride *int64
	// getNilContentLength forces the response ContentLength to be nil,
	// regardless of getContentLengthOverride.
	getNilContentLength bool

	deleteInput *s3.DeleteObjectInput
	deleteErr   error

	// real-s3 fallback for Presign (its concrete request.Request is hard to fake)
	realS3 *s3.S3
}

func (f *fakeS3) PutObjectWithContext(_ aws.Context, in *s3.PutObjectInput, _ ...request.Option) (*s3.PutObjectOutput, error) {
	f.putInput = in
	if f.putErr != nil {
		return &s3.PutObjectOutput{}, f.putErr
	}
	return &s3.PutObjectOutput{}, nil
}

func (f *fakeS3) GetObjectWithContext(_ aws.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
	f.getInput = in
	if f.getErr != nil {
		return nil, f.getErr
	}
	body := io.ReadCloser(io.NopCloser(bytes.NewReader(f.getBody)))
	if f.getBodyOverride != nil {
		body = f.getBodyOverride
	}
	size := int64(len(f.getBody))
	out := &s3.GetObjectOutput{Body: body}
	switch {
	case f.getNilContentLength:
		out.ContentLength = nil
	case f.getContentLengthOverride != nil:
		out.ContentLength = f.getContentLengthOverride
	default:
		out.ContentLength = &size
	}
	return out, nil
}

func (f *fakeS3) DeleteObjectWithContext(_ aws.Context, in *s3.DeleteObjectInput, _ ...request.Option) (*s3.DeleteObjectOutput, error) {
	f.deleteInput = in
	if f.deleteErr != nil {
		return &s3.DeleteObjectOutput{}, f.deleteErr
	}
	return &s3.DeleteObjectOutput{}, nil
}

func (f *fakeS3) GetObjectRequest(in *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	return f.realS3.GetObjectRequest(in)
}

func sampleConfig() Config {
	return Config{AWSS3: AWSS3Config{
		Region:          "ap-southeast-1",
		BucketName:      "test-bucket",
		AccessKeyID:     "AKIA-test",
		SecretAccessKey: "secret-test",
		PresignDuration: 5 * time.Minute,
	}}
}

// realS3ForTests builds a *s3.S3 against the (offline) credentials in the test
// config — Presign signs locally, no network call is made.
func realS3ForTests(t *testing.T, cfg Config) *s3.S3 {
	t.Helper()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSS3.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AWSS3.AccessKeyID, cfg.AWSS3.SecretAccessKey, ""),
	})
	require.NoError(t, err)
	return s3.New(sess)
}

func newStorage(t *testing.T, fake *fakeS3) *storage {
	t.Helper()
	cfg := sampleConfig()
	fake.realS3 = realS3ForTests(t, cfg)
	return &storage{
		s3:     fake,
		config: cfg,
		log:    logger.Init(logger.Config{}),
	}
}

// ---------------- Upload ---------------- //

func TestUpload_Success(t *testing.T) {
	fake := &fakeS3{}
	s := newStorage(t, fake)

	got, err := s.Upload(context.Background(), "k", "file.json", "application/json", []byte("hello"))
	require.NoError(t, err)
	assert.Contains(t, got, "test-bucket")
	assert.Contains(t, got, "ap-southeast-1")
	assert.True(t, strings.HasSuffix(got, "/k"))

	require.NotNil(t, fake.putInput)
	assert.Equal(t, "test-bucket", aws.StringValue(fake.putInput.Bucket))
	assert.Equal(t, "k", aws.StringValue(fake.putInput.Key))
	assert.Equal(t, "application/json", aws.StringValue(fake.putInput.ContentType))
}

func TestUpload_Error(t *testing.T) {
	fake := &fakeS3{putErr: errors.New("s3 down")}
	s := newStorage(t, fake)
	_, err := s.Upload(context.Background(), "k", "f", "txt", []byte("x"))
	assert.Error(t, err)
}

// ---------------- Download ---------------- //

func TestDownload_Success(t *testing.T) {
	body := []byte("payload-data")
	fake := &fakeS3{getBody: body}
	s := newStorage(t, fake)

	got, err := s.Download(context.Background(), "any-key")
	require.NoError(t, err)
	assert.Equal(t, body, got)
	require.NotNil(t, fake.getInput)
	assert.Equal(t, "any-key", aws.StringValue(fake.getInput.Key))
}

func TestDownload_Error(t *testing.T) {
	fake := &fakeS3{getErr: errors.New("not found")}
	s := newStorage(t, fake)
	_, err := s.Download(context.Background(), "k")
	assert.Error(t, err)
}

// erroringReader returns a slice of bytes once and then a non-EOF error,
// simulating a network/body failure partway through the stream.
type erroringReader struct {
	chunk   []byte
	sent    bool
	failErr error
}

func (e *erroringReader) Read(p []byte) (int, error) {
	if !e.sent {
		e.sent = true
		n := copy(p, e.chunk)
		return n, nil
	}
	return 0, e.failErr
}

func (e *erroringReader) Close() error { return nil }

func TestDownload_PropagatesReadError(t *testing.T) {
	fake := &fakeS3{
		getBodyOverride: &erroringReader{chunk: []byte("partial"), failErr: errors.New("network glitch")},
	}
	// Force ContentLength to a non-matching positive value to make sure
	// the body-read path is exercised even with a non-nil length.
	cl := int64(1024)
	fake.getContentLengthOverride = &cl
	s := newStorage(t, fake)

	_, err := s.Download(context.Background(), "k")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network glitch")
}

func TestDownload_NilContentLengthIsSafe(t *testing.T) {
	body := []byte("payload")
	fake := &fakeS3{
		getBodyOverride:     io.NopCloser(bytes.NewReader(body)),
		getNilContentLength: true,
	}
	s := newStorage(t, fake)

	got, err := s.Download(context.Background(), "k")
	require.NoError(t, err)
	assert.Equal(t, body, got)
}

// ---------------- Delete ---------------- //

func TestDelete_Success(t *testing.T) {
	fake := &fakeS3{}
	s := newStorage(t, fake)

	require.NoError(t, s.Delete(context.Background(), "k"))
	require.NotNil(t, fake.deleteInput)
	assert.Equal(t, "k", aws.StringValue(fake.deleteInput.Key))
	assert.Equal(t, "test-bucket", aws.StringValue(fake.deleteInput.Bucket))
}

func TestDelete_Error(t *testing.T) {
	fake := &fakeS3{deleteErr: errors.New("denied")}
	s := newStorage(t, fake)
	assert.Error(t, s.Delete(context.Background(), "k"))
}

// ---------------- Presign (uses real *s3.S3 — no network) ---------------- //

func TestGetPresignedUrl_Default(t *testing.T) {
	fake := &fakeS3{}
	s := newStorage(t, fake)

	url, err := s.GetPresignedUrl(context.Background(), "my/key.png")
	require.NoError(t, err)
	assert.Contains(t, url, "test-bucket")
	assert.Contains(t, url, "X-Amz-Signature=") // standard presign query param
}

func TestGetPresignedUrlWithDuration(t *testing.T) {
	fake := &fakeS3{}
	s := newStorage(t, fake)

	url, err := s.GetPresignedUrlWithDuration(context.Background(), "k", 10*time.Minute)
	require.NoError(t, err)
	assert.Contains(t, url, "X-Amz-Signature=")
}

func TestGetPresignedUrl_InvalidDurationReturnsError(t *testing.T) {
	fake := &fakeS3{}
	s := newStorage(t, fake)

	// Override the default duration with an invalid (zero) one via the
	// dedicated method since Presign rejects expire<=0.
	_, err := s.GetPresignedUrlWithDuration(context.Background(), "k", 0)
	assert.Error(t, err)
}

// ---------------- Init ---------------- //

func TestInit_BuildsRealS3(t *testing.T) {
	cfg := sampleConfig()
	got := Init(cfg, logger.Init(logger.Config{}))
	require.NotNil(t, got)
}
