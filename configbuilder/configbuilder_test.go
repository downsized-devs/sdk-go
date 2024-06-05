//go:build integration
// +build integration

package configbuilder

import (
	"log"
	"os"
	"testing"

	"github.com/hlubek/readercomp"
	"github.com/stretchr/testify/assert"
)

const (
	configfile   string = "./conf.json"
	templatefile string = "./files/conf.template"
	appnamespace string = "aquaheroservice"
)

func Test_configBuilder_BuildConfig(t *testing.T) {
	tests := []struct {
		name       string
		opt        Options
		wantErr    bool
		actualFile string
	}{
		{
			name: "success",
			opt: Options{
				Env:          "", // fill these in if you want to run the test
				Key:          "", // fill these in if you want to run the test
				Secret:       "", // fill these in if you want to run the test
				Region:       "", // fill these in if you want to run the test
				TemplateFile: templatefile,
				ConfigFile:   configfile,
				Namespace:    appnamespace,
			},
			wantErr:    false,
			actualFile: "./files/actual_conf.json",
		},
	}
	for _, tt := range tests {
		defer func() {
			e := os.Remove(tt.opt.ConfigFile)
			if e != nil {
				log.Fatal(e)
			}
		}()
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.opt).BuildConfig()
			if tt.wantErr {
				assert.Equal(t, assert.FileExists(t, configfile), false)
			}
			result, err := readercomp.FilesEqual(configfile, tt.actualFile)
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, result, true)
		})
	}
}
