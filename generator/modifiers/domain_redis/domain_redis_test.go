package domainRedis

import (
	"os"
	"testing"
)

func Test_domainRedis_Replace(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll("templates", os.ModePerm)
			if err != nil {
				t.Errorf("domainRedis.Replace() error = %v", err)
			}

			file, err := os.Create("templates/domain_redis.tmpl")
			if err != nil {
				t.Errorf("domainRedis.Replace() error = %v", err)
			}
			file.Close()
			d := &domainRedis{
				EntityName:          tt.fields.EntityName,
				EntityNameUpper:     tt.fields.EntityNameUpper,
				EntityNameSnakeCase: tt.fields.EntityNameSnakeCase,
				Location:            tt.fields.Location,
			}
			if err := d.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("domainRedis.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat("generate/src/business/domain/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + "_redis.go"); err == nil {
				if err := os.RemoveAll("generate/src/business/domain/" + tt.fields.EntityNameSnakeCase + "/" + tt.fields.EntityNameSnakeCase + "_redis.go"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/domain_redis.tmpl"); err == nil {
				if err := os.RemoveAll("templates/domain_redis.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}
