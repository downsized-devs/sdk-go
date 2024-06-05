package usecaseItem

import (
	"os"
	"testing"
)

func Test_usecaseItem_Replace(t *testing.T) {
	type fields struct {
		EntityName          string
		EntityNameUpper     string
		EntityNameSnakeCase string
		Location            string
		Api                 []string
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
				Location:            "generate",
				Api: []string{
					"create",
					"get",
					"delete",
					"activate",
					"edit",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		err := os.MkdirAll("templates/usecase_item", os.ModePerm)
		if err != nil {
			t.Errorf("usecaseItem.Replace() error = %v", err)
		}

		file, err := os.Create("templates/usecase_item/usecase_item.tmpl")
		if err != nil {
			t.Errorf("usecaseItem.Replace() error = %v", err)
		}
		file.Close()
		t.Run(tt.name, func(t *testing.T) {
			d := &usecaseItem{
				EntityName:          tt.fields.EntityName,
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Location:            tt.fields.Location,
				Api:                 tt.fields.Api,
			}
			if err := d.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("usecaseItem.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(d.Location + "/src/business/usecase/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + ".go"); err == nil {
				if err := os.RemoveAll(d.Location + "/src/business/usecase/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + ".go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/usecase_item/usecase_item.tmpl"); err == nil {
				if err := os.RemoveAll("templates/usecase_item/usecase_item.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}

func Test_usecaseItem_AppendInterfaceAndFunction(t *testing.T) {
	type fields struct {
		EntityName          string
		EntityNameUpper     string
		EntityNameSnakeCase string
		Location            string
		Api                 []string
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
				Location:            "generate",
				Api: []string{
					"create",
					"get",
					"delete",
					"activate",
					"edit",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.fields.Location+"/src/business/usecase/"+tt.fields.EntityNameSnakeCase+"/", os.ModePerm)
			if err != nil {
				t.Errorf("usecaseItem.AppendInterfaceAndFunction() error = %v", err)
			}

			file, err := os.Create(tt.fields.Location + "/src/business/usecase/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + ".go")
			if err != nil {
				t.Errorf("usecaseItem.AppendInterfaceAndFunction() error = %v", err)
			}
			file.WriteString("type Interface interface {")
			file.Close()
			for _, api := range tt.fields.Api {
				err := os.MkdirAll("templates/usecase_item/", os.ModePerm)
				if err != nil {
					t.Errorf("secaseItem.AppendInterfaceAndFunction() error = %v", err)
				}
				file, err := os.Create("templates/usecase_item/usecase_item_" + api + "_function.tmpl")
				if err != nil {
					t.Errorf("secaseItem.AppendInterfaceAndFunction() error = %v", err)
				}
				file.Close()
			}
			d := &usecaseItem{
				EntityName:          tt.fields.EntityName,
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Location:            tt.fields.Location,
				Api:                 tt.fields.Api,
			}
			if err := d.AppendInterfaceAndFunction(); (err != nil) != tt.wantErr {
				t.Errorf("usecaseItem.AppendInterfaceAndFunction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.fields.Location + "/src/business/usecase/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + ".go"); err == nil {
				if err := os.RemoveAll(tt.fields.Location + "/src/business/usecase/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + ".go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			for _, api := range tt.fields.Api {
				if _, err := os.Stat("templates/usecase_item/usecase_item_" + api + "_function.tmpl"); err == nil {
					if err := os.RemoveAll("templates/usecase_item/usecase_item_" + api + "_function.tmpl"); err != nil {
						t.Errorf("failed to remove item")
					}
				}
			}
		})
	}
}
