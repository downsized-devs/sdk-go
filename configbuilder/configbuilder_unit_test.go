package configbuilder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit-test counterpart to configbuilder_test.go (which is build-tagged and
// uses a fixture template). These exercise the same code paths using
// per-test temp directories so they run by default.

func TestInit_PanicsOnEmptyEnv(t *testing.T) {
	assert.Panics(t, func() {
		Init(Options{Env: ""})
	})
}

func TestInit_OK(t *testing.T) {
	b := Init(Options{Env: "test"})
	assert.NotNil(t, b)
}

func TestBuildConfig_TemplateMissing(t *testing.T) {
	b := Init(Options{
		Env:          "test",
		TemplateFile: "/this/does/not/exist.template",
		ConfigFile:   filepath.Join(t.TempDir(), "out.json"),
	})
	assert.Panics(t, b.BuildConfig)
}

func TestBuildConfig_OK(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "in.template")
	out := filepath.Join(dir, "out.json")

	// Simple template with no variables — mustache renders unchanged.
	if err := os.WriteFile(tpl, []byte(`{"key":"value"}`), 0o600); err != nil {
		t.Fatalf("write template: %v", err)
	}

	b := Init(Options{Env: "test", TemplateFile: tpl, ConfigFile: out})
	b.BuildConfig()

	got, err := os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, `{"key":"value"}`, string(got))
}

func TestBuildConfig_MustacheError(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "bad.template")
	out := filepath.Join(dir, "out.json")

	// Unclosed mustache tag triggers a render error which BuildConfig
	// converts into a panic.
	if err := os.WriteFile(tpl, []byte(`{{unclosed`), 0o600); err != nil {
		t.Fatalf("write template: %v", err)
	}

	b := Init(Options{Env: "test", TemplateFile: tpl, ConfigFile: out})
	assert.Panics(t, b.BuildConfig)
}

func TestBuildConfig_CreateOutputError(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "in.template")
	if err := os.WriteFile(tpl, []byte(`{"k":"v"}`), 0o600); err != nil {
		t.Fatalf("write template: %v", err)
	}
	// Output path lives in a directory that doesn't exist — os.Create fails.
	out := filepath.Join(dir, "missing-dir", "out.json")
	b := Init(Options{Env: "test", TemplateFile: tpl, ConfigFile: out})
	assert.Panics(t, b.BuildConfig)
}

func TestBuildConfig_UnknownVariableIsBlank(t *testing.T) {
	// Without an explicit viper.BindEnv, AutomaticEnv does not populate
	// AllSettings(), so unbound template variables render as empty.
	dir := t.TempDir()
	tpl := filepath.Join(dir, "in.template")
	out := filepath.Join(dir, "out.json")

	if err := os.WriteFile(tpl, []byte(`{"version":"{{some_var}}"}`), 0o600); err != nil {
		t.Fatalf("write template: %v", err)
	}

	b := Init(Options{Env: "test", TemplateFile: tpl, ConfigFile: out})
	b.BuildConfig()

	got, err := os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, `{"version":""}`, string(got))
}
