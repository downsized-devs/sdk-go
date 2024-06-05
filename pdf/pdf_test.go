package pdf

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/downsized-devs/sdk-go/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func Test_pdf_SetPasswordFile(t *testing.T) {
	type fields struct {
		config *model.Configuration
		log    log.Interface
	}

	pdfFile := "example.pdf"
	pdfBytes, err := os.ReadFile(pdfFile)
	if err != nil {
		fmt.Println("Error reading PDF file:", err)
		return
	}

	type args struct {
		ctx      context.Context
		password string
		data     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				config: model.NewDefaultConfiguration(),
				log:    nil,
			},
			args: args{
				ctx:      context.Background(),
				password: "password",
				data:     pdfBytes,
			},
			want:    pdfBytes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pdf{
				config: tt.fields.config,
				log:    tt.fields.log,
			}
			_, err := p.SetPasswordFile(tt.args.ctx, tt.args.password, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("pdf.SetPasswordFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
