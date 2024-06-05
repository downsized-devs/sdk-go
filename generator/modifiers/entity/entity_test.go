package entity

import (
	"os"
	"testing"
)

func Test_entity_Replace(t *testing.T) {
	type fields struct {
		EntityNameUpper     string
		EntityNameSnakeCase string
		EntityName          string
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
			t.Errorf("entity.Replace() error = %v", err)
		}

		file, err := os.Create("templates/entity.tmpl")
		if err != nil {
			t.Errorf("entity.Replace() error = %v", err)
		}
		file.Close()
		t.Run(tt.name, func(t *testing.T) {
			e := &entity{
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				EntityName:          tt.fields.EntityName,
				Location:            tt.fields.Location,
			}
			if err := e.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("entity.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat("generate/src/business/entity/" + tt.fields.EntityNameSnakeCase + ".go"); err == nil {
				if err := os.RemoveAll("generate/src/business/entity/" + tt.fields.EntityNameSnakeCase + ".go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/entity.tmpl"); err == nil {
				if err := os.RemoveAll("templates/entity.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}
