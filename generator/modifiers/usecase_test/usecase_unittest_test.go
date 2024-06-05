package usecaseUnittest

import (
	"os"
	"testing"
)

func Test_usecaseTest_Replace(t *testing.T) {
	type fields struct {
		EntityName          string
		EntityNameUpper     string
		EntityNameSnakeCase string
		Api                 []string
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
				EntityName:          "name",
				EntityNameUpper:     "Name",
				EntityNameSnakeCase: "name",
				Api: []string{
					"create",
					"get",
					"delete",
					"activate",
					"edit",
				},
				Location: "generate",
			},
		},
	}
	for _, tt := range tests {
		err := os.MkdirAll("templates/usecase_test", os.ModePerm)
		if err != nil {
			t.Errorf("usecaseTest.Replace() error = %v", err)
		}

		file, err := os.Create("templates/usecase_test/usecase_test.tmpl")
		if err != nil {
			t.Errorf("usecaseTest.Replace() error = %v", err)
		}
		file.Close()
		t.Run(tt.name, func(t *testing.T) {
			d := &usecaseTest{
				EntityName:          tt.fields.EntityName,
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Api:                 tt.fields.Api,
				Location:            tt.fields.Location,
			}
			if err := d.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("usecaseTest.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(d.Location + "/src/business/usecase/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + "_test.go"); err == nil {
				if err := os.RemoveAll(d.Location + "/src/business/usecase/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + "_test.go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/usecase_test/usecase_test.tmpl"); err == nil {
				if err := os.RemoveAll("templates/usecase_test/usecase_test.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}

func Test_usecaseTest_AppendInterfaceAndFunction(t *testing.T) {
	type fields struct {
		EntityName          string
		EntityNameUpper     string
		EntityNameSnakeCase string
		Api                 []string
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
				EntityName:          "name",
				EntityNameUpper:     "Name",
				EntityNameSnakeCase: "name",
				Api: []string{
					"create",
					"get",
					"delete",
					"activate",
					"edit",
				},
				Location: "generate",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.fields.Location+"/src/business/usecase/"+tt.fields.EntityNameSnakeCase+"/", os.ModePerm)
			if err != nil {
				t.Errorf("usecaseTest.AppendInterfaceAndFunction() error = %v", err)
			}

			file, err := os.Create(tt.fields.Location + "/src/business/usecase/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + "_test.go")
			if err != nil {
				t.Errorf("usecaseTest.AppendInterfaceAndFunction() error = %v", err)
			}
			file.Close()
			for _, api := range tt.fields.Api {
				err := os.MkdirAll("templates/usecase_test/", os.ModePerm)
				if err != nil {
					t.Errorf("usecaseTest.AppendInterfaceAndFunction() error = %v", err)
				}
				file, err := os.Create("templates/usecase_test/usecase_test_" + api + "_function.tmpl")
				if err != nil {
					t.Errorf("usecaseTest.AppendInterfaceAndFunction() error = %v", err)
				}
				file.Close()
			}
			d := &usecaseTest{
				EntityName:          tt.fields.EntityName,
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Api:                 tt.fields.Api,
				Location:            tt.fields.Location,
			}
			if err := d.AppendInterfaceAndFunction(); (err != nil) != tt.wantErr {
				t.Errorf("usecaseTest.AppendInterfaceAndFunction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.fields.Location + "/src/business/usecase/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + "_test.go"); err == nil {
				if err := os.RemoveAll(tt.fields.Location + "/src/business/usecase/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + "_test.go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			for _, api := range tt.fields.Api {
				if _, err := os.Stat("templates/usecase_test/usecase_test_" + api + "_function.tmpl"); err == nil {
					if err := os.RemoveAll("templates/usecase_test/usecase_test_" + api + "_function.tmpl"); err != nil {
						t.Errorf("failed to remove item")
					}
				}
			}
		})
	}
}
