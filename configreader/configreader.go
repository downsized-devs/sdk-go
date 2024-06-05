package configreader

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/downsized-devs/sdk-go/files"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	JSONType string = "json"
	YAMLType string = "yaml"
)

type Interface interface {
	ReadConfig(cfg interface{})
	AllSettings() map[string]interface{}
}

type Options struct {
	// location of the main config file
	ConfigFile string
	// additional configuration to append to the main config
	AdditionalConfig []AdditionalConfigOptions
}
type AdditionalConfigOptions struct {
	// key is the location in the main config to append the additional config
	// with dot delimiter. e.g. 'Parser.ExcelOptions'
	ConfigKey string
	// location of the additional config
	ConfigFile string
}

type configReader struct {
	viper *viper.Viper
	opt   Options
}

func Init(opt Options) Interface {
	vp := viper.New()
	vp.SetConfigFile(opt.ConfigFile)
	if err := vp.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error found during reading file. err: %w", err))
	}

	for _, addConf := range opt.AdditionalConfig {
		vpAddConf := viper.New()
		vpAddConf.SetConfigFile(addConf.ConfigFile)
		if err := vpAddConf.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error found during reading additional config file. err: %w", err))
		}
		if addConf.ConfigKey == "" {
			vp.MergeConfigMap(vpAddConf.AllSettings())
		} else {
			vp.Set(strings.ToLower(addConf.ConfigKey), vpAddConf.AllSettings())
		}
	}

	c := &configReader{
		viper: vp,
		opt:   opt,
	}

	return c
}

func (c *configReader) mergeEnvConfig() {
	var svcVersion string
	switch {
	case os.Getenv("AQUAHERO_SERVICE_VERSION") != "":
		svcVersion = os.Getenv("AQUAHERO_SERVICE_VERSION")
	case os.Getenv("SERVICE_VERSION") != "":
		svcVersion = os.Getenv("SERVICE_VERSION")
	default:
		svcVersion = "dev"
	}
	meta := c.viper.GetStringMap("meta")
	meta["version"] = svcVersion
	c.viper.Set("meta", meta)
}

func (c *configReader) resolveJSONRef() {
	refmap := make(map[string]interface{})
	refregxp := regexp.MustCompile("^\\$ref:#\\/(.*)$")
	for _, k := range c.viper.AllKeys() {
		refpath := c.viper.GetString(k)
		if refregxp.MatchString(refpath) {
			v, ok := refmap[refpath]
			if !ok {
				refkey := refregxp.ReplaceAllString(refpath, "$1")
				refkey = strings.ToLower(strings.ReplaceAll(refkey, "/", "."))
				refmap[refpath] = c.viper.Get(refkey)
				c.viper.Set(k, refmap[refpath])
			} else {
				c.viper.Set(k, v)
			}
		}
	}
}

func (c *configReader) ReadConfig(cfg interface{}) {
	c.mergeEnvConfig()

	if files.GetExtension(filepath.Base(c.opt.ConfigFile)) == JSONType {
		c.resolveJSONRef()
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig(&cfg))
	if err != nil {
		panic(fmt.Errorf("fatal error found during decoder creation. err: %w", err))
	}

	if err := decoder.Decode(c.viper.AllSettings()); err != nil {
		panic(fmt.Errorf("fatal error found during unmarshaling config. err: %w", err))
	}
}

func (c *configReader) AllSettings() map[string]interface{} {
	return c.viper.AllSettings()
}

// internally modified string to duration parser hooks function to handle empty string
func stringToTimeDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Duration(5)) {
			return data, nil
		}

		// If data is empty string return zero duration
		if data.(string) == "" {
			return time.Duration(0), nil
		}

		// Convert it by parsing
		return time.ParseDuration(data.(string))
	}
}

// Modified default decoder config to avoid errors parsing duration from empty string
// This decoder is using internally modified string to duration parser hooks function
func decoderConfig(output interface{}, opts ...viper.DecoderConfigOption) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
