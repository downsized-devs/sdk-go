package restinItialize

import (
	"os"
	"testing"
)

func Test_restInitialize_Initialize(t *testing.T) {
	type fields struct {
		EntityNameUpper      string
		EntityNameSnakeCase  string
		EntityNameLowerSpace string
		EntityNameLowerDash  string
		Location             string
		Api                  []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				EntityNameUpper:      "Name",
				EntityNameSnakeCase:  "name",
				EntityNameLowerSpace: "name",
				EntityNameLowerDash:  "name",
				Location:             "generate",
				Api: []string{
					"create",
					"get",
					"delete",
					"activate",
					"edit",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.fields.Location+"/src/handler/rest/", os.ModePerm)
			if err != nil {
				t.Errorf("handlerRest.AppendInterfaceAndFunction() error = %v", err)
			}

			file, err := os.Create(tt.fields.Location + "/src/handler/rest/rest.go")
			if err != nil {
				t.Errorf("handlerRest.AppendInterfaceAndFunction() error = %v", err)
			}
			file.WriteString("func (r *rest) Register() {\n\n}")
			file.Close()
			d := &restInitialize{
				EntityNameUpper:      tt.fields.EntityNameUpper,
				EntityNameSnakeCase:  tt.fields.EntityNameSnakeCase,
				EntityNameLowerSpace: tt.fields.EntityNameLowerSpace,
				EntityNameLowerDash:  tt.fields.EntityNameLowerDash,
				Location:             tt.fields.Location,
				Api:                  tt.fields.Api,
			}
			if err := d.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("restInitialize.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.fields.Location + "/src/handler/rest/rest.go"); err == nil {
				if err := os.RemoveAll("generate/src/handler/rest/rest.go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}
