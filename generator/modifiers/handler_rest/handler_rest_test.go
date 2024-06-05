package handlerRest

import (
	"os"
	"testing"
)

func Test_handlerRest_Replace(t *testing.T) {
	type fields struct {
		EntityName           string
		EntityNameUpper      string
		EntityNameSnakeCase  string
		EntityNameLowerSpace string
		EntityNameUpperSpace string
		EntityNameLowerDash  string
		HandlerParamItems    []string
		RouterItems          string
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
				EntityName:           "name",
				EntityNameUpper:      "Name",
				EntityNameSnakeCase:  "name",
				EntityNameLowerSpace: "name",
				EntityNameUpperSpace: "Name",
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll("templates/handler_rest", os.ModePerm)
			if err != nil {
				t.Errorf("handlerRest.Replace() error = %v", err)
			}

			file, err := os.Create("templates/handler_rest/handler_rest.tmpl")
			if err != nil {
				t.Errorf("handlerRest.Replace() error = %v", err)
			}
			file.Close()
			d := &handlerRest{
				EntityName:           tt.fields.EntityName,
				EntityNameUpper:      tt.fields.EntityNameUpper,
				EntityNameSnakeCase:  tt.fields.EntityNameSnakeCase,
				EntityNameLowerSpace: tt.fields.EntityNameLowerSpace,
				EntityNameUpperSpace: tt.fields.EntityNameUpperSpace,
				EntityNameLowerDash:  tt.fields.EntityNameLowerDash,
				HandlerParamItems:    tt.fields.HandlerParamItems,
				RouterItems:          tt.fields.RouterItems,
				Location:             tt.fields.Location,
				Api:                  tt.fields.Api,
			}
			if err := d.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("handlerRest.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat("generate/src/handler/rest/name.go"); err == nil {
				if err := os.RemoveAll("generate/src/handler/rest/name.go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/handler_rest/handler_rest.tmpl"); err == nil {
				if err := os.RemoveAll("templates/handler_rest/handler_rest.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}

func Test_handlerRest_AppendInterfaceAndFunction(t *testing.T) {
	type fields struct {
		EntityName           string
		EntityNameUpper      string
		EntityNameSnakeCase  string
		EntityNameLowerSpace string
		EntityNameUpperSpace string
		EntityNameLowerDash  string
		HandlerParamItems    []string
		RouterItems          string
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
				EntityName:           "name",
				EntityNameUpper:      "Name",
				EntityNameSnakeCase:  "name",
				EntityNameLowerSpace: "name",
				EntityNameUpperSpace: "Name",
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.fields.Location+"/src/handler/rest/", os.ModePerm)
			if err != nil {
				t.Errorf("handlerRest.AppendInterfaceAndFunction() error = %v", err)
			}

			file, err := os.Create(tt.fields.Location + "/src/handler/rest/" + tt.fields.EntityNameSnakeCase + ".go")
			if err != nil {
				t.Errorf("handlerRest.AppendInterfaceAndFunction() error = %v", err)
			}
			file.Close()
			for _, api := range tt.fields.Api {
				err := os.MkdirAll("templates/handler_rest/", os.ModePerm)
				if err != nil {
					t.Errorf("handlerRest.AppendInterfaceAndFunction() error = %v", err)
				}
				file, err := os.Create("templates/handler_rest/handler_rest_" + api + "_function.tmpl")
				if err != nil {
					t.Errorf("handlerRest.AppendInterfaceAndFunction() error = %v", err)
				}
				file.Close()
			}
			d := &handlerRest{
				EntityName:           tt.fields.EntityName,
				EntityNameUpper:      tt.fields.EntityNameUpper,
				EntityNameSnakeCase:  tt.fields.EntityNameSnakeCase,
				EntityNameLowerSpace: tt.fields.EntityNameLowerSpace,
				EntityNameUpperSpace: tt.fields.EntityNameUpperSpace,
				EntityNameLowerDash:  tt.fields.EntityNameLowerDash,
				HandlerParamItems:    tt.fields.HandlerParamItems,
				RouterItems:          tt.fields.RouterItems,
				Location:             tt.fields.Location,
				Api:                  tt.fields.Api,
			}
			if err := d.AppendInterfaceAndFunction(); (err != nil) != tt.wantErr {
				t.Errorf("handlerRest.AppendInterfaceAndFunction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat("generate/src/handler/rest/" + tt.fields.EntityNameSnakeCase + ".go"); err == nil {
				if err := os.RemoveAll("generate/src/handler/rest/" + tt.fields.EntityNameSnakeCase + ".go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			for _, api := range tt.fields.Api {
				if _, err := os.Stat("templates/handler_rest/handler_rest_" + api + "_function.tmpl"); err == nil {
					if err := os.RemoveAll("templates/handler_rest/handler_rest_" + api + "_function.tmpl"); err != nil {
						t.Errorf("failed to remove item")
					}
				}
			}
		})
	}
}
