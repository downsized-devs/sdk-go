package domainQuery

import (
	"os"
	"text/template"
)

type domainQuery struct {
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
	return &domainQuery{
		EntityName:          param.EntityName,
		EntityNameUpper:     param.EntityNameUpper,
		EntityNameSnakeCase: param.EntityNameSnakeCase,
		Location:            param.Location,
	}
}

func (d *domainQuery) Replace() error {
	paths := []string{
		"templates/domain_query.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(d.Location+"/src/business/domain/"+d.EntityNameSnakeCase, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(d.Location + "/src/business/domain/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + "_query.go")
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
