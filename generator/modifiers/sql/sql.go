package sql

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

type sql struct {
	EntityNameUpperSpace string
	EntityNameSnakeCase  string
	Location             string
}

type Interface interface {
	Replace() error
}

type Params struct {
	EntityNameUpperSpace string
	EntityNameSnakeCase  string
	Location             string
}

func Init(param Params) Interface {
	return &sql{
		EntityNameUpperSpace: param.EntityNameUpperSpace,
		EntityNameSnakeCase:  param.EntityNameSnakeCase,
		Location:             param.Location,
	}
}

func (s *sql) Replace() error {
	paths := []string{
		"templates/sql.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(s.Location+"/docs/sql/", os.ModePerm)
	if err != nil {
		return err
	}

	sqlName := fmt.Sprintf("%d%02d%02d", time.Now().Year()%1e2, time.Now().Month(), time.Now().Day())

	count := 0
	for {
		tempSqlName := sqlName + "_" + fmt.Sprintf("%02d", count)
		if _, err := os.Stat(s.Location + "/docs/sql/" + tempSqlName + ".sql"); err != nil {
			sqlName = tempSqlName
			break
		}
		count++
	}

	file, err := os.Create(s.Location + "/docs/sql/" + sqlName + ".sql")
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Execute(file, s)
	if err != nil {
		return err
	}
	return nil
}
