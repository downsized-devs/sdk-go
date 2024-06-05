package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/downsized-devs/sdk-go/log"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

const (
	outputFileName = "output.pdf"
)

type PdfInterface interface {
	SetPasswordFile(ctx context.Context, password string, data []byte) ([]byte, error)
}

type pdf struct {
	config *model.Configuration
	log    log.Interface
}

func Init(log log.Interface) PdfInterface {
	p := &pdf{
		config: model.NewDefaultConfiguration(),
		log:    log,
	}

	return p
}

func (p *pdf) SetPasswordFile(ctx context.Context, password string, data []byte) ([]byte, error) {
	p.config.UserPW = password
	p.config.OwnerPW = password

	// Setting encryption mode to AES-256
	p.config.EncryptUsingAES = true
	p.config.EncryptKeyLength = 256

	// Write the byte array to a file
	err := os.WriteFile(outputFileName, data, 0777)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("Failed to write PDF file: %v", err.Error()))
		return data, err
	}

	// remove file after process encrypted
	defer func() {
		err = os.Remove(outputFileName)
		if err != nil {
			p.log.Error(ctx, err)
			return
		}
	}()

	// Encrypt the file
	err = api.EncryptFile(outputFileName, outputFileName, p.config)
	if err != nil {
		return data, err
	}

	res := bytes.NewBuffer(nil)
	file, err := os.Open(outputFileName)
	if err != nil {
		return data, err
	}

	if _, err := io.Copy(res, file); err != nil {
		return data, err
	}

	return res.Bytes(), nil
}
