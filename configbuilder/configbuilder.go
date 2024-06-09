package configbuilder

import (
	"os"

	"github.com/cbroglie/mustache"
	"github.com/downsized-devs/sdk-go/files"
	"github.com/spf13/viper"
)

type Options struct {
	Namespace              string
	Env                    string
	Key                    string
	Secret                 string
	Region                 string
	TemplateFile           string
	ConfigFile             string
	IgnoreEmptyConfigError bool
}

type Interface interface {
	BuildConfig()
}

type configBuilder struct {
	opt Options
}

func Init(opt Options) Interface {
	if opt.Env == "" {
		panic("Environment variable is not set")
	}

	return &configBuilder{
		opt: opt,
	}
}

func (b *configBuilder) BuildConfig() {
	if !files.IsExist(b.opt.TemplateFile) {
		panic("Template file not found")
	}

	params := viper.New()

	// TODO: set up config source

	body, err := os.ReadFile(b.opt.TemplateFile)
	if err != nil {
		panic(err.Error())
	}

	conf, err := mustache.Render(string(body), true, params.AllSettings())
	if err != nil {
		panic(err.Error())
	}

	f, err := os.Create(b.opt.ConfigFile)
	if err != nil {
		panic(err.Error())
	}

	_, err = f.Write([]byte(conf))
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}
}
