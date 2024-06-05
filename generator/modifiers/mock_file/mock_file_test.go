package mockfile

import (
	"os"
	"testing"
)

func Test_mockFile_Replace(t *testing.T) {
	type fields struct {
		EntityName          string
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
				EntityName:          "turtlePond",
				EntityNameUpper:     "TurtlePond",
				EntityNameSnakeCase: "turtle_pond",
				Location:            "generate",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		err := os.MkdirAll("templates", os.ModePerm)
		if err != nil {
			t.Errorf("mockFile.Replace() error = %v", err)
		}

		file, err := os.Create("templates/mock_file.tmpl")
		if err != nil {
			t.Errorf("mockFile.Replace() error = %v", err)
		}
		file.Close()
		t.Run(tt.name, func(t *testing.T) {
			d := &mockFile{
				EntityName:          tt.fields.EntityName,
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Location:            tt.fields.Location,
			}
			if err := d.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("mockFile.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat("generate/src/business/domain/mock/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + ".go"); err == nil {
				if err := os.RemoveAll("generate/src/business/domain/mock/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + ".go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/mock_file.tmpl"); err == nil {
				if err := os.RemoveAll("templates/mock_file.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}

func Test_mockFile_Initialize(t *testing.T) {
	type fields struct {
		EntityName          string
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
				EntityName:          "turtlePond",
				EntityNameUpper:     "TurtlePond",
				EntityNameSnakeCase: "turtle_pond",
				Location:            "generate",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		err := os.MkdirAll(tt.fields.Location+"/", os.ModePerm)
		if err != nil {
			t.Errorf("mockFile.Initialize() error = %v", err)
		}

		file, err := os.Create(tt.fields.Location + "/Makefile")
		if err != nil {
			t.Errorf("mockFile.Initialize() error = %v", err)
		}
		_, err = file.WriteString("mock-all:")
		if err != nil {
			t.Errorf("mockFile.Initialize() error = %v", err)
		}
		file.Close()

		t.Run(tt.name, func(t *testing.T) {
			d := &mockFile{
				EntityName:          tt.fields.EntityName,
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Location:            tt.fields.Location,
			}
			if err := d.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("mockFile.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		if _, err := os.Stat(tt.fields.Location + "/Makefile"); err == nil {
			if err := os.RemoveAll(tt.fields.Location + "/Makefile"); err != nil {
				t.Errorf("failed to remove item")
			}
		}

	}
}
