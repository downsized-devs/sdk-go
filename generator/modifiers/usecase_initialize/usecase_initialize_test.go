package usecaseInitialize

import (
	"os"
	"testing"
)

func Test_usecaseInitialize_Initialize(t *testing.T) {
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
			name: "success",
			fields: fields{
				EntityNameUpper:     "Name",
				EntityNameSnakeCase: "name",
				Location:            "generate",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.fields.Location+"/src/business/usecase", os.ModePerm)
			if err != nil {
				t.Errorf("usecaseInitialize.Initialize error = %v", err)
			}

			file, err := os.Create(tt.fields.Location + "/src/business/usecase/usecase.go")
			if err != nil {
				t.Errorf("usecaseInitialize.Initialize error = %v", err)
			}
			file.WriteString("import (\n\ntype Usecases struct {\n\nuc := &Usecases{")
			file.Close()
			d := &usecaseInitialize{
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Location:            tt.fields.Location,
			}
			if err := d.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("usecaseInitialize.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.fields.Location + "/src/business/usecase/usecase.go"); err == nil {
				if err := os.RemoveAll(tt.fields.Location + "/src/business/usecase/usecase.go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}
