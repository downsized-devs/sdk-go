package domainRedis

import (
	"os"
	"text/template"
)

type domainRedis struct {
	EntityName          string
	EntityNameUpper     string
	EntityNameSnakeCase string
	Location            string
}

type Interface interface {
	Replace() error
}

type Params struct {
	EntityName          string
	EntityNameUpper     string
	EntityNameSnakeCase string
	Location            string
}

func Init(param Params) Interface {
	return &domainRedis{
		EntityName:          param.EntityName,
		EntityNameUpper:     param.EntityNameUpper,
		EntityNameSnakeCase: param.EntityNameSnakeCase,
		Location:            param.Location,
	}
}

func (d *domainRedis) Replace() error {
	paths := []string{
		"templates/domain_redis.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(d.Location+"/src/business/domain/"+d.EntityNameSnakeCase, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(d.Location + "/src/business/domain/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + "_redis.go")
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Execute(file, d)
	if err != nil {
		return err
	}
	return nil
}
