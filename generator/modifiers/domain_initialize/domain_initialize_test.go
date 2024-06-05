package domaininItialize

import (
	"os"
	"testing"
)

func Test_domainInitialize_Initialize(t *testing.T) {
	type fields struct {
		EntityNameUpper     string
		EntityNameSnakeCase string
		Location            string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "no file",
			fields: fields{
				EntityNameUpper:     "TurtlePond",
				EntityNameSnakeCase: "turtle_pond",
				Location:            "generate",
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				EntityNameUpper:     "TurtlePond",
				EntityNameSnakeCase: "turtle_pond",
				Location:            "test_folder",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success" {
				err := os.MkdirAll(tt.fields.Location+"/src/business/domain", os.ModePerm)
				if err != nil {
					t.Errorf("domainInitialize.Initialize error = %v", err)
				}

				file, err := os.Create(tt.fields.Location + "/src/business/domain/domain.go")
				if err != nil {
					t.Errorf("domainInitialize.Initialize error = %v", err)
				}
				file.WriteString("import (\n\ntype Domains struct {\n\nd := &Domains{")
				file.Close()
			}
			d := &domainInitialize{
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Location:            tt.fields.Location,
			}
			if err := d.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("domainInitialize.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.fields.Location + "/src/business/domain/domain.go"); err == nil {
				if err := os.RemoveAll(tt.fields.Location + "/src/business/domain/domain.go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}
