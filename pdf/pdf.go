// Package pdf is a thin in-memory wrapper around pdfcpu (encryption, merge,
// split, watermark, page count) plus ledongthuc/pdf for text extraction.
//
// All methods take and return []byte; no temp files leak out except where
// pdfcpu itself only exposes a directory-based API (Split), in which case a
// temp directory is created and removed before returning.
package pdf

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
	ledongpdf "github.com/ledongthuc/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Default watermark description string passed through to pdfcpu. Callers who
// need a different font/size/opacity can use AddTextWatermarkWithDesc.
const defaultWatermarkDesc = "font:Helvetica, points:24, opacity:0.5, scale:0.5"

type Interface interface {
	Encrypt(ctx context.Context, data []byte, password string) ([]byte, error)
	RemovePassword(ctx context.Context, data []byte, password string) ([]byte, error)
	Merge(ctx context.Context, parts ...[]byte) ([]byte, error)
	Split(ctx context.Context, data []byte, span int) ([][]byte, error)
	AddTextWatermark(ctx context.Context, data []byte, text string) ([]byte, error)
	ExtractText(ctx context.Context, data []byte) (string, error)
	PageCount(ctx context.Context, data []byte) (int, error)
}

type pdf struct {
	log logger.Interface
}

func Init(log logger.Interface) Interface {
	return &pdf{log: log}
}

func (p *pdf) Encrypt(_ context.Context, data []byte, password string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: empty input")
	}
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password
	conf.EncryptUsingAES = true
	conf.EncryptKeyLength = 256

	out := &bytes.Buffer{}
	if err := api.Encrypt(bytes.NewReader(data), out, conf); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (p *pdf) RemovePassword(_ context.Context, data []byte, password string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: empty input")
	}
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password

	out := &bytes.Buffer{}
	if err := api.Decrypt(bytes.NewReader(data), out, conf); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (p *pdf) Merge(_ context.Context, parts ...[]byte) ([]byte, error) {
	if len(parts) == 0 {
		return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: no parts to merge")
	}
	if len(parts) == 1 {
		// Nothing to merge — return a copy so callers can rely on a non-aliased buffer.
		dup := make([]byte, len(parts[0]))
		copy(dup, parts[0])
		return dup, nil
	}

	rscs := make([]io.ReadSeeker, len(parts))
	for i, b := range parts {
		if len(b) == 0 {
			return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: merge part %d is empty", i)
		}
		rscs[i] = bytes.NewReader(b)
	}
	out := &bytes.Buffer{}
	if err := api.MergeRaw(rscs, out, false, model.NewDefaultConfiguration()); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (p *pdf) Split(_ context.Context, data []byte, span int) ([][]byte, error) {
	if span < 1 {
		return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: span must be >= 1")
	}
	if len(data) == 0 {
		return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: empty input")
	}

	dir, err := os.MkdirTemp("", "sdkgo-pdf-split-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	inPath := filepath.Join(dir, "in.pdf")
	if err := os.WriteFile(inPath, data, 0o600); err != nil {
		return nil, err
	}
	if err := api.SplitFile(inPath, dir, span, model.NewDefaultConfiguration()); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, e := range entries {
		if e.IsDir() || e.Name() == "in.pdf" {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names) // pdfcpu emits in.pdf_1.pdf, in.pdf_2.pdf, ... so lexical = page order

	chunks := make([][]byte, 0, len(names))
	for _, n := range names {
		b, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, b)
	}
	return chunks, nil
}

func (p *pdf) AddTextWatermark(_ context.Context, data []byte, text string) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: empty input")
	}
	if text == "" {
		return nil, errors.NewWithCode(codes.CodeBadRequest, "pdf: empty watermark text")
	}

	wm, err := api.TextWatermark(text, defaultWatermarkDesc, true, false, types.POINTS)
	if err != nil {
		return nil, err
	}
	out := &bytes.Buffer{}
	if err := api.AddWatermarks(bytes.NewReader(data), out, nil, wm, model.NewDefaultConfiguration()); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (p *pdf) ExtractText(_ context.Context, data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.NewWithCode(codes.CodeBadRequest, "pdf: empty input")
	}

	r, err := ledongpdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	total := r.NumPage()
	for i := 1; i <= total; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		txt, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		buf.WriteString(txt)
	}
	return buf.String(), nil
}

func (p *pdf) PageCount(_ context.Context, data []byte) (int, error) {
	if len(data) == 0 {
		return 0, errors.NewWithCode(codes.CodeBadRequest, "pdf: empty input")
	}
	return api.PageCount(bytes.NewReader(data), model.NewDefaultConfiguration())
}

var _ Interface = (*pdf)(nil)
