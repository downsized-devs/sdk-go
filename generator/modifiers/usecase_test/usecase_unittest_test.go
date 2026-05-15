package usecaseUnittest

import (
	"os"
	"path/filepath"
	"strings"
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

func Test_usecaseTest_AppendInterfaceAndFunction_ReadOpenError(t *testing.T) {
	d := &usecaseTest{
		EntityNameSnakeCase: "missing_entity",
		Location:            "does/not/exist",
	}
	err := d.AppendInterfaceAndFunction()
	if err == nil {
		t.Fatal("expected error when target file is missing during read phase, got nil")
	}
	if !os.IsNotExist(err) && !strings.Contains(err.Error(), "no such file") {
		t.Errorf("expected file-not-found error, got %v", err)
	}
}

func Test_usecaseTest_AppendInterfaceAndFunction_DirectoryAsTarget(t *testing.T) {
	tmpDir := t.TempDir()
	entityDir := filepath.Join(tmpDir, "src", "business", "usecase", "fake", "fake_test.go")
	if err := os.MkdirAll(entityDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	d := &usecaseTest{
		EntityNameSnakeCase: "fake",
		Location:            tmpDir,
	}
	if err := d.AppendInterfaceAndFunction(); err == nil {
		t.Fatal("expected error when target path is a directory, got nil")
	}
}

func Test_usecaseTest_AppendInterfaceAndFunction_EmptyAPI(t *testing.T) {
	tmpDir := t.TempDir()
	entityDir := filepath.Join(tmpDir, "src", "business", "usecase", "entity")
	if err := os.MkdirAll(entityDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	target := filepath.Join(entityDir, "entity_test.go")
	if err := os.WriteFile(target, []byte("package entity_test\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	d := &usecaseTest{
		EntityNameSnakeCase: "entity",
		Location:            tmpDir,
		Api:                 nil,
	}
	if err := d.AppendInterfaceAndFunction(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("file was truncated, want non-empty content")
	}
}

func Test_usecaseTest_getFunction(t *testing.T) {
	tmpDir := t.TempDir()
	d := &usecaseTest{}

	t.Run("returns lines on success", func(t *testing.T) {
		path := filepath.Join(tmpDir, "ok.txt")
		if err := os.WriteFile(path, []byte("one\ntwo\n"), 0o644); err != nil {
			t.Fatalf("setup: %v", err)
		}
		lines, err := d.getFunction(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 2 || lines[0] != "one" || lines[1] != "two" {
			t.Errorf("unexpected lines: %v", lines)
		}
	})

	t.Run("returns error when file is missing", func(t *testing.T) {
		_, err := d.getFunction(filepath.Join(tmpDir, "nope.txt"))
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("returns error when path is a directory", func(t *testing.T) {
		dirPath := filepath.Join(tmpDir, "asdir")
		if err := os.Mkdir(dirPath, 0o755); err != nil {
			t.Fatalf("setup: %v", err)
		}
		if _, err := d.getFunction(dirPath); err == nil {
			t.Fatal("expected error when opening a directory, got nil")
		}
	})
}
