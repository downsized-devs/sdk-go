package handlerRest

import (
	"bufio"
	"os"
	"strings"
	"text/template"
)

type handlerRest struct {
	EntityName           string
	EntityNameUpper      string
	EntityNameSnakeCase  string
	EntityNameLowerSpace string
	EntityNameUpperSpace string
	EntityNameLowerDash  string
	HandlerParamItems    []string
	RouterItems          string
	Location             string
	Api                  []string
}

type Interface interface {
	Replace() error
	AppendInterfaceAndFunction() error
}

type Params struct {
	EntityName           string
	EntityNameUpper      string
	EntityNameSnakeCase  string
	EntityNameLowerSpace string
	EntityNameUpperSpace string
	EntityNameLowerDash  string
	Location             string
	Api                  []string
}

func Init(param Params) Interface {
	return &handlerRest{
		EntityName:           param.EntityName,
		EntityNameUpper:      param.EntityNameUpper,
		EntityNameSnakeCase:  param.EntityNameSnakeCase,
		EntityNameLowerSpace: param.EntityNameLowerSpace,
		EntityNameUpperSpace: param.EntityNameUpperSpace,
		EntityNameLowerDash:  param.EntityNameLowerDash,
		Location:             param.Location,
		Api:                  param.Api,
	}
}

func (d *handlerRest) Replace() error {
	paths := []string{
		"templates/handler_rest/handler_rest.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(d.Location+"/src/handler/rest/", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(d.Location + "/src/handler/rest/" + d.EntityNameSnakeCase + ".go")
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

func (d *handlerRest) AppendInterfaceAndFunction() error {
	f, err := os.OpenFile(d.Location+"/src/handler/rest/"+d.EntityNameSnakeCase+".go", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		ln := scanner.Text()
		lines = append(lines, ln)
	}
	f.Close()
	for _, api := range d.Api {
		apiLower := strings.ToLower(api)
		apiLower = strings.ReplaceAll(apiLower, " ", "")
		lines, err = d.appendFunction(lines, apiLower, "templates/handler_rest/handler_rest_temp.tmpl")
		if err != nil {
			return err
		}
	}
	content := strings.Join(lines, "\n")
	f, err = os.OpenFile(d.Location+"/src/handler/rest/"+d.EntityNameSnakeCase+".go", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteAt([]byte(content), 0)
	if err != nil {
		return err
	}
	return nil
}

func (d *handlerRest) replaceFunction(templatePath, generatedPath string) (*os.File, error) {
	paths := []string{
		templatePath,
	}
	t := template.Must(template.ParseFiles(paths...))
	file, err := os.Create(generatedPath)
	if err != nil {
		return nil, err
	}
	err = t.Execute(file, d)
	if err != nil {
		return nil, err
	}
	file.Close()
	return file, nil
}

func (d *handlerRest) getFunction(path string) ([]string, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		ln := scanner.Text()
		lines = append(lines, ln)
	}
	f.Close()
	return lines, nil
}

func (d *handlerRest) appendFunction(lines []string, api, generatedPath string) ([]string, error) {
	result := lines
	file, err := d.replaceFunction("templates/handler_rest/handler_rest_"+api+"_function.tmpl", generatedPath)
	if err != nil {
		return result, err
	}
	tempFunction, err := d.getFunction(generatedPath)
	if err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := scanner.Text()
		tempFunction = append(tempFunction, ln)
	}
	result = append(result, tempFunction...)
	file.Close()
	err = os.Remove(generatedPath)
	if err != nil {
		return result, err
	}
	return result, nil
}
