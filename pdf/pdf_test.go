package pdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
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

// buildPDF assembles a tiny PDF with the given object bodies. Each body must be
// the full "N 0 obj ... endobj\n" string. Returns the bytes plus byte offsets
// for the supplied objects (1-indexed). It also writes a valid xref table and
// startxref/EOF so ledongpdf.NewReader accepts it.
func buildPDF(bodies ...string) []byte {
	var buf bytes.Buffer
	write := func(s string) int {
		off := buf.Len()
		buf.WriteString(s)
		return off
	}

	write("%PDF-1.4\n")
	offsets := make([]int, len(bodies))
	for i, body := range bodies {
		offsets[i] = write(body)
	}
	xref := buf.Len()
	write(fmt.Sprintf("xref\n0 %d\n", len(bodies)+1))
	write("0000000000 65535 f \n")
	for _, off := range offsets {
		write(fmt.Sprintf("%010d 00000 n \n", off))
	}
	write(fmt.Sprintf("trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF",
		len(bodies)+1, xref))
	return buf.Bytes()
}

// nullPagePDF reports Count=1 but has an empty /Kids array, so r.Page(1) falls
// through Reader.Page's search loop and returns Page{} — i.e. p.V.IsNull().
func nullPagePDF() []byte {
	return buildPDF(
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
		"2 0 obj\n<< /Type /Pages /Kids [] /Count 1 >>\nendobj\n",
	)
}

// brokenContentPDF has a real /Page whose /Contents stream contains a bare
// "Tj" operator with no operands. GetPlainText's Interpret callback panics on
// `len(args) != 1`, and the deferred recover surfaces that as an error.
func brokenContentPDF() []byte {
	const content = "Tj"
	return buildPDF(
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
		"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /Contents 4 0 R /MediaBox [0 0 612 792] >>\nendobj\n",
		fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(content), content),
	)
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

// TestSplit_MkdirTempError forces os.MkdirTemp to fail by pointing TMPDIR at a
// path that doesn't exist. This exercises the first error return in Split.
func TestSplit_MkdirTempError(t *testing.T) {
	t.Setenv("TMPDIR", "/nonexistent/path/that/should/never/exist/sdkgo-pdf-test")
	p := newPDF()
	_, err := p.Split(context.Background(), loadExample(t), 1)
	assert.Error(t, err)
}

// TestSplit_WriteFileError swaps the WriteFile hook to drive Split's
// "could not stage input" error branch — unreachable under a fresh temp dir
// otherwise.
func TestSplit_WriteFileError(t *testing.T) {
	sentinel := errors.New("write failed")
	orig := osWriteFile
	osWriteFile = func(string, []byte, fs.FileMode) error { return sentinel }
	t.Cleanup(func() { osWriteFile = orig })

	p := newPDF()
	_, err := p.Split(context.Background(), loadExample(t), 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, sentinel)
}

// TestSplit_ReadDirError swaps the ReadDir hook so Split fails after a
// successful split. The temp dir is cleaned up via the production defer.
func TestSplit_ReadDirError(t *testing.T) {
	sentinel := errors.New("readdir failed")
	orig := osReadDir
	osReadDir = func(string) ([]fs.DirEntry, error) { return nil, sentinel }
	t.Cleanup(func() { osReadDir = orig })

	p := newPDF()
	_, err := p.Split(context.Background(), loadExample(t), 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, sentinel)
}

// TestSplit_ReadFileError swaps ReadFile so the chunk-collection loop in Split
// fails on the first chunk emitted by pdfcpu.
func TestSplit_ReadFileError(t *testing.T) {
	sentinel := errors.New("readfile failed")
	orig := osReadFile
	osReadFile = func(string) ([]byte, error) { return nil, sentinel }
	t.Cleanup(func() { osReadFile = orig })

	p := newPDF()
	_, err := p.Split(context.Background(), loadExample(t), 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, sentinel)
}

func TestAddTextWatermark(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	stamped, err := p.AddTextWatermark(context.Background(), data, "CONFIDENTIAL")
	require.NoError(t, err)
	assert.NotEmpty(t, stamped)
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

// TestAddTextWatermark_DescriptorParseError swaps the package-level watermark
// description with one pdfcpu's parser will reject ("font" without a value),
// covering the api.TextWatermark error branch that the hardcoded valid
// descriptor never exercises.
func TestAddTextWatermark_DescriptorParseError(t *testing.T) {
	orig := defaultWatermarkDesc
	defaultWatermarkDesc = "not-a-valid-descriptor-format"
	t.Cleanup(func() { defaultWatermarkDesc = orig })

	p := newPDF()
	_, err := p.AddTextWatermark(context.Background(), loadExample(t), "WM")
	assert.Error(t, err)
}

func TestExtractText(t *testing.T) {
	p := newPDF()
	data := loadExample(t)

	txt, err := p.ExtractText(context.Background(), data)
	require.NoError(t, err)
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

// TestExtractText_NullPageSkipped feeds in a PDF whose /Pages dictionary
// advertises a Count of 1 but has an empty /Kids list. ledongpdf.NumPage
// returns 1, but r.Page(1).V.IsNull() is true, so the loop hits `continue`
// and ExtractText returns an empty string with no error.
func TestExtractText_NullPageSkipped(t *testing.T) {
	p := newPDF()
	txt, err := p.ExtractText(context.Background(), nullPagePDF())
	require.NoError(t, err)
	assert.Equal(t, "", txt)
}

// TestExtractText_GetPlainTextError feeds a PDF with a real page whose
// /Contents stream is a bare "Tj" operator. GetPlainText's interpreter panics
// on the empty operand stack and the deferred recover surfaces it as an error.
func TestExtractText_GetPlainTextError(t *testing.T) {
	p := newPDF()
	_, err := p.ExtractText(context.Background(), brokenContentPDF())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Tj")
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
