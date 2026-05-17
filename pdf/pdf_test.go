package pdf

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadExample returns the bytes of example.pdf shipped in the test data.
func loadExample(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("example.pdf")
	require.NoError(t, err)
	return b
}

func newPDF() Interface {
	return Init(logger.Init(logger.Config{}))
}

func TestInit_ReturnsInterface(t *testing.T) {
	p := newPDF()
	assert.NotNil(t, p)
}

func TestPageCount(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	n, err := p.PageCount(context.Background(), data)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, n, 1)
}

func TestPageCount_EmptyInput(t *testing.T) {
	p := newPDF()
	_, err := p.PageCount(context.Background(), nil)
	assert.Error(t, err)
}

func TestEncrypt_RoundTrip(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	encrypted, err := p.Encrypt(context.Background(), data, "secret")
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.False(t, bytes.Equal(data, encrypted), "encrypted output should differ from input")

	decrypted, err := p.RemovePassword(context.Background(), encrypted, "secret")
	require.NoError(t, err)
	assert.NotEmpty(t, decrypted)

	// After decrypt → re-encrypt cycle we should be able to read the page count again.
	n, err := p.PageCount(context.Background(), decrypted)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, n, 1)
}

func TestEncrypt_EmptyInput(t *testing.T) {
	p := newPDF()
	_, err := p.Encrypt(context.Background(), nil, "x")
	assert.Error(t, err)
}

func TestRemovePassword_EmptyInput(t *testing.T) {
	p := newPDF()
	_, err := p.RemovePassword(context.Background(), nil, "x")
	assert.Error(t, err)
}

func TestMerge_TwoCopies(t *testing.T) {
	p := newPDF()
	data := loadExample(t)
	originalPages, err := p.PageCount(context.Background(), data)
	require.NoError(t, err)

	merged, err := p.Merge(context.Background(), data, data)
	require.NoError(t, err)
	assert.NotEmpty(t, merged)

	got, err := p.PageCount(context.Background(), merged)
	require.NoError(t, err)
	assert.Equal(t, originalPages*2, got, "merging a PDF with itself doubles the page count")
}

func TestMerge_SinglePartReturnsCopy(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	got, err := p.Merge(context.Background(), data)
	require.NoError(t, err)
	assert.Equal(t, data, got)

	// Mutating the result must not mutate the input.
	if len(got) > 0 {
		got[0] ^= 0xFF
		assert.NotEqual(t, got[0], data[0])
	}
}

func TestMerge_NoParts(t *testing.T) {
	p := newPDF()
	_, err := p.Merge(context.Background())
	assert.Error(t, err)
}

func TestMerge_EmptyPart(t *testing.T) {
	p := newPDF()
	data := loadExample(t)
	_, err := p.Merge(context.Background(), data, nil)
	assert.Error(t, err)
}

func TestSplit(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	pages, err := p.PageCount(context.Background(), data)
	require.NoError(t, err)
	if pages < 2 {
		// The fixture is single-page; build a multi-page fixture by merging.
		data, err = p.Merge(context.Background(), data, data, data)
		require.NoError(t, err)
		pages = 3
	}

	chunks, err := p.Split(context.Background(), data, 1)
	require.NoError(t, err)
	assert.Equal(t, pages, len(chunks), "span=1 produces one chunk per page")

	for i, c := range chunks {
		n, err := p.PageCount(context.Background(), c)
		require.NoError(t, err, "chunk %d", i)
		assert.Equal(t, 1, n)
	}
}

func TestSplit_InvalidSpan(t *testing.T) {
	p := newPDF()
	data := loadExample(t)
	_, err := p.Split(context.Background(), data, 0)
	assert.Error(t, err)
}

func TestSplit_EmptyInput(t *testing.T) {
	p := newPDF()
	_, err := p.Split(context.Background(), nil, 1)
	assert.Error(t, err)
}

func TestAddTextWatermark(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	stamped, err := p.AddTextWatermark(context.Background(), data, "CONFIDENTIAL")
	require.NoError(t, err)
	assert.NotEmpty(t, stamped)
	// Watermarking changes the byte stream.
	assert.False(t, bytes.Equal(data, stamped))
}

func TestAddTextWatermark_EmptyText(t *testing.T) {
	p := newPDF()
	data := loadExample(t)
	_, err := p.AddTextWatermark(context.Background(), data, "")
	assert.Error(t, err)
}

func TestAddTextWatermark_EmptyInput(t *testing.T) {
	p := newPDF()
	_, err := p.AddTextWatermark(context.Background(), nil, "x")
	assert.Error(t, err)
}

func TestExtractText(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	txt, err := p.ExtractText(context.Background(), data)
	require.NoError(t, err)
	// ExtractText output for the fixture isn't strictly defined; assert basic invariants only.
	_ = txt
	assert.NotPanics(t, func() {
		_ = strings.TrimSpace(txt)
	})
}

func TestExtractText_EmptyInput(t *testing.T) {
	p := newPDF()
	_, err := p.ExtractText(context.Background(), nil)
	assert.Error(t, err)
}

// garbage is a non-empty byte slice that is syntactically not a PDF. Used to
// drive the api.* and ledongpdf.NewReader failure paths past the len==0 guards.
var garbage = []byte("not a pdf, just noise to defeat the empty-input guard")

func TestEncrypt_InvalidPDF(t *testing.T) {
	p := newPDF()
	_, err := p.Encrypt(context.Background(), garbage, "secret")
	assert.Error(t, err)
}

func TestRemovePassword_InvalidPDF(t *testing.T) {
	p := newPDF()
	_, err := p.RemovePassword(context.Background(), garbage, "secret")
	assert.Error(t, err)
}

func TestMerge_InvalidPart(t *testing.T) {
	p := newPDF()
	data := loadExample(t)
	// Two non-empty inputs forces api.MergeRaw, where the garbage part fails.
	_, err := p.Merge(context.Background(), data, garbage)
	assert.Error(t, err)
}

func TestSplit_InvalidPDF(t *testing.T) {
	p := newPDF()
	_, err := p.Split(context.Background(), garbage, 1)
	assert.Error(t, err)
}

func TestAddTextWatermark_InvalidPDF(t *testing.T) {
	p := newPDF()
	_, err := p.AddTextWatermark(context.Background(), garbage, "WM")
	assert.Error(t, err)
}

func TestExtractText_InvalidPDF(t *testing.T) {
	p := newPDF()
	_, err := p.ExtractText(context.Background(), garbage)
	assert.Error(t, err)
}

func TestPageCount_InvalidPDF(t *testing.T) {
	p := newPDF()
	_, err := p.PageCount(context.Background(), garbage)
	assert.Error(t, err)
}
