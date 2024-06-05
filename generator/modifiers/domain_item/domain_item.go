package domainItem

import (
	"os"
	"text/template"
)

type domainItem struct {
	EntityName           string
	EntityNameUpper      string
	EntityNameSnakeCase  string
	EntityNameLowerSpace string
	Location             string
}

type Interface interface {
	Replace() error
}

type Params struct {
	EntityName           string
	EntityNameUpper      string
	EntityNameSnakeCase  string
	EntityNameLowerSpace string
	Location             string
}

func Init(param Params) Interface {
	return &domainItem{
		EntityName:           param.EntityName,
		EntityNameUpper:      param.EntityNameUpper,
		EntityNameSnakeCase:  param.EntityNameSnakeCase,
		EntityNameLowerSpace: param.EntityNameLowerSpace,
		Location:             param.Location,
	}
}

func (d *domainItem) Replace() error {
	paths := []string{
		"templates/domain_item.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(d.Location+"/src/business/domain/"+d.EntityNameSnakeCase, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(d.Location + "/src/business/domain/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + ".go")
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
