package configbuilder

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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
	ssm *ssm.SSM
}

func Init(opt Options) Interface {
	var sess *session.Session

	if opt.Env == "" {
		panic("Environment variable is not set")
	}

	if opt.Region == "" {
		opt.Region = "ap-southeast-1"
	}

	// if access key and secret not found in env, get credentials from local or metadata
	// this behaviour is intended to support all local, staging, and production environment
	// session credentials behaviour: https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
	if opt.Key != "" && opt.Secret != "" {
		sess = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(opt.Region),
			Credentials: credentials.NewStaticCredentials(opt.Key, opt.Secret, ""),
		}))
	} else {
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(opt.Region),
		}))
	}

	ssm := ssm.New(sess)
	return &configBuilder{
		opt: opt,
		ssm: ssm,
	}
}

func (b *configBuilder) BuildConfig() {
	if !files.IsExist(b.opt.TemplateFile) {
		panic("Template file not found")
	}

	// Expire current credentials to refresh session in case previous session is still active
	// This is to avoid getting cached value on parameter store
	b.ssm.Config.Credentials.Expire()

	ssmparams := []*ssm.Parameter{}
	ssmres, err := b.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(fmt.Sprintf("/service/%s/", b.opt.Env)),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		panic(err.Error())
	}
	ssmparams = append(ssmparams, ssmres.Parameters...)
	for ssmres.NextToken != nil {
		ssmres, err = b.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
			Path:           aws.String(fmt.Sprintf("/service/%s/", b.opt.Env)),
			NextToken:      ssmres.NextToken,
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			panic(err.Error())
		}
		ssmparams = append(ssmparams, ssmres.Parameters...)
	}

	if len(ssmparams) < 1 && !b.opt.IgnoreEmptyConfigError {
		panic("No configuration found. Please check your environment variable")
	}

	params := viper.New()
	// Set Current Namespace Params
	for _, p := range ssmparams {
		val := *p.Value
		keyregxp := regexp.MustCompile(fmt.Sprintf("^\\/service\\/%s\\/%s\\/(.*)$", b.opt.Env, b.opt.Namespace))
		key := strings.ReplaceAll(keyregxp.ReplaceAllString(*p.Name, "$1"), "/", ".")
		params.Set(key, val)
	}

	// Set Global Namespace Params
	for _, p := range ssmparams {
		val := *p.Value
		key := strings.ReplaceAll(*p.Name, "/", ".")
		key = strings.ReplaceAll(key, fmt.Sprintf(".service.%s.", b.opt.Env), "service.")
		params.Set(key, val)
	}

	body, err := ioutil.ReadFile(b.opt.TemplateFile)
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
