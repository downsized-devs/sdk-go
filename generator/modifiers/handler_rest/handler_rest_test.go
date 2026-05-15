package handlerRest

import (
	"os"
	"path/filepath"
	"strings"
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

func Test_handlerRest_AppendInterfaceAndFunction_ReadOpenError(t *testing.T) {
	d := &handlerRest{
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

func Test_handlerRest_AppendInterfaceAndFunction_DirectoryAsTarget(t *testing.T) {
	// OpenFile against a directory in O_RDWR mode fails on darwin/linux,
	// exercising the early-return branch of the read phase.
	tmpDir := t.TempDir()
	entityDir := filepath.Join(tmpDir, "src", "handler", "rest", "fake.go")
	if err := os.MkdirAll(entityDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	d := &handlerRest{
		EntityNameSnakeCase: "fake",
		Location:            tmpDir,
	}
	if err := d.AppendInterfaceAndFunction(); err == nil {
		t.Fatal("expected error when target path is a directory, got nil")
	}
}

func Test_handlerRest_AppendInterfaceAndFunction_EmptyAPI(t *testing.T) {
	// Empty Api skips the append loop, so the read phase, write phase, and
	// deferred close all run to completion on a real file. Exercises the
	// happy path of the named-return + defer-close pattern.
	tmpDir := t.TempDir()
	restDir := filepath.Join(tmpDir, "src", "handler", "rest")
	if err := os.MkdirAll(restDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	target := filepath.Join(restDir, "entity.go")
	original := []byte("package rest\n// existing\n")
	if err := os.WriteFile(target, original, 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	d := &handlerRest{
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
	// strings.Join with no separator additions still preserves original
	// content; main concern is that the file isn't truncated or corrupted.
	if len(got) == 0 {
		t.Errorf("file was truncated, want non-empty content")
	}
}

func Test_handlerRest_getFunction(t *testing.T) {
	tmpDir := t.TempDir()
	d := &handlerRest{}

	t.Run("returns lines on success", func(t *testing.T) {
		path := filepath.Join(tmpDir, "ok.txt")
		if err := os.WriteFile(path, []byte("alpha\nbeta\ngamma\n"), 0o644); err != nil {
			t.Fatalf("setup: %v", err)
		}
		lines, err := d.getFunction(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"alpha", "beta", "gamma"}
		if len(lines) != len(want) {
			t.Fatalf("got %d lines, want %d (%v)", len(lines), len(want), lines)
		}
		for i, ln := range lines {
			if ln != want[i] {
				t.Errorf("line %d = %q, want %q", i, ln, want[i])
			}
		}
	})

	t.Run("returns error when file is missing", func(t *testing.T) {
		_, err := d.getFunction(filepath.Join(tmpDir, "does-not-exist.txt"))
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("returns error when path is a directory", func(t *testing.T) {
		dirPath := filepath.Join(tmpDir, "as-dir")
		if err := os.Mkdir(dirPath, 0o755); err != nil {
			t.Fatalf("setup: %v", err)
		}
		_, err := d.getFunction(dirPath)
		if err == nil {
			t.Fatal("expected error when opening a directory, got nil")
		}
	})
}
