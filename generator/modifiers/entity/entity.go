package entity

import (
	"os"
	"text/template"
)

type entity struct {
	EntityNameUpper     string
	EntityNameSnakeCase string
	EntityName          string
	Location            string
}

type Interface interface {
	Replace() error
}

type Params struct {
	EntityNameUpper     string
	EntityNameSnakeCase string
	EntityName          string
	Location            string
}

func Init(param Params) Interface {
	return &entity{
		EntityNameUpper:     param.EntityNameUpper,
		EntityNameSnakeCase: param.EntityNameSnakeCase,
		EntityName:          param.EntityName,
		Location:            param.Location,
	}
}

func (e *entity) Replace() error {
	paths := []string{
		"templates/entity.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(e.Location+"/src/business/entity/", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(e.Location + "/src/business/entity/" + e.EntityNameSnakeCase + ".go")
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Execute(file, e)
	if err != nil {
		return err
	}
	return nil
}
