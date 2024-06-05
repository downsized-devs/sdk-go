package sql

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_sql_Replace(t *testing.T) {
	type fields struct {
		EntityNameUpperSpace string
		EntityNameSnakeCase  string
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
				EntityNameUpperSpace: "Name",
				EntityNameSnakeCase:  "name",
				Location:             "generate",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.fields.Location+"/docs/sql/", os.ModePerm)
			if err != nil {
				t.Errorf("sql.Replace() error = %v", err)
			}

			file, err := os.Create(tt.fields.Location + "/docs/sql/" + fmt.Sprintf("%d%02d%02d", time.Now().Year()%1e2, time.Now().Month(), time.Now().Day()) + "_00.sql")
			if err != nil {
				t.Errorf("sql.Replace() error = %v", err)
			}
			file.Close()
			err = os.MkdirAll("templates", os.ModePerm)
			if err != nil {
				t.Errorf("usecaseTest.Replace() error = %v", err)
			}

			file, err = os.Create("templates/sql.tmpl")
			if err != nil {
				t.Errorf("usecaseTest.Replace() error = %v", err)
			}
			file.Close()
			s := &sql{
				EntityNameUpperSpace: tt.fields.EntityNameUpperSpace,
				EntityNameSnakeCase:  tt.fields.EntityNameSnakeCase,
				Location:             tt.fields.Location,
			}
			if err := s.Replace(); (err != nil) != tt.wantErr {
				t.Errorf("sql.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.fields.Location + "/docs/sql/" + fmt.Sprintf("%d%02d%02d", time.Now().Year()%1e2, time.Now().Month(), time.Now().Day()) + "_00.sql"); err == nil {
				if err := os.RemoveAll(tt.fields.Location + "/docs/sql/" + fmt.Sprintf("%d%02d%02d", time.Now().Year()%1e2, time.Now().Month(), time.Now().Day()) + "_00.sql"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat(tt.fields.Location + "/docs/sql/" + fmt.Sprintf("%d%02d%02d", time.Now().Year()%1e2, time.Now().Month(), time.Now().Day()) + "_01.sql"); err == nil {
				if err := os.RemoveAll(tt.fields.Location + "/docs/sql/" + fmt.Sprintf("%d%02d%02d", time.Now().Year()%1e2, time.Now().Month(), time.Now().Day()) + "_01.sql"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
			if _, err := os.Stat("templates/sql.tmpl"); err == nil {
				if err := os.RemoveAll("templates/sql.tmpl"); err != nil {
					t.Errorf("failed to remove item")
				}
			}
		})
	}
}
