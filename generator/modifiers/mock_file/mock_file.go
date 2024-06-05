package mockfile

import (
	"bufio"
	"os"
	"strings"
	"text/template"
)

type mockFile struct {
	EntityName          string
	EntityNameUpper     string
	EntityNameSnakeCase string

	Location string
}

type Interface interface {
	Replace() error
	Initialize() error
}

type Params struct {
	EntityName          string
	EntityNameUpper     string
	EntityNameSnakeCase string

	Location string
}

func Init(param Params) Interface {
	return &mockFile{
		EntityName:          param.EntityName,
		EntityNameUpper:     param.EntityNameUpper,
		EntityNameSnakeCase: param.EntityNameSnakeCase,
		Location:            param.Location,
	}
}

func (d *mockFile) Replace() error {
	paths := []string{
		"templates/mock_file.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(d.Location+"/src/business/domain/mock/"+d.EntityNameSnakeCase, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(d.Location + "/src/business/domain/mock/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + ".go")
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

func (d *mockFile) Initialize() error {
	f, err := os.OpenFile(d.Location+"/Makefile", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		ln := scanner.Text()
		lines = append(lines, ln)
		if strings.Contains(ln, "mock-all:") {
			lines = append(lines, "	@make mock domain="+d.EntityNameSnakeCase)
		}
	}

	content := strings.Join(lines, "\n")
	_, err = f.WriteAt([]byte(content), 0)
	if err != nil {
		return err
	}
	return nil
}
