package domainItem

import (
	"os"
	"testing"
)

func Test_domainItem_Replace(t *testing.T) {
	type fields struct {
		EntityName           string
		EntityNameUpper      string
		EntityNameSnakeCase  string
		EntityNameLowerSpace string
		Location             string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				EntityName:           "turtlePond",
				EntityNameUpper:      "TurtlePond",
				EntityNameSnakeCase:  "turtle_pond",
				EntityNameLowerSpace: "turtle pond",
				Location:             "generate",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll("templates", os.ModePerm)
			if err != nil {
				t.Errorf("domainItem.Replace() error = %v", err)
			}

			file, err := os.Create("templates/domain_item.tmpl")
			if err != nil {
				t.Errorf("domainItem.Replace() error = %v", err)
			}
			file.Close()
			d := &domainItem{
				EntityName:           tt.fields.EntityName,
				EntityNameUpper:      tt.fields.EntityNameUpper,
				EntityNameSnakeCase:  tt.fields.EntityNameSnakeCase,
				EntityNameLowerSpace: tt.fields.EntityNameLowerSpace,
				Location:             tt.fields.Location,
			}
			if err := d.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("domainItem.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat("generate/src/business/domain/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + ".go"); err == nil {
				if err := os.RemoveAll("generate/src/business/domain/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + ".go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/domain_item.tmpl"); err == nil {
				if err := os.RemoveAll("templates/domain_item.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}
