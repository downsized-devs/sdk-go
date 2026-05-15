package usecaseUnittest

import (
	"bufio"
	"os"
	"strings"
	"text/template"
)

type usecaseTest struct {
	EntityName          string
	EntityNameUpper     string
	EntityNameSnakeCase string
	Api                 []string
	Location            string
}

type Interface interface {
	Replace() error
	AppendInterfaceAndFunction() error
}

type Params struct {
	EntityName          string
	EntityNameUpper     string
	EntityNameSnakeCase string
	Api                 []string
	Location            string
}

func Init(param Params) Interface {
	return &usecaseTest{
		EntityName:          param.EntityName,
		EntityNameUpper:     param.EntityNameUpper,
		EntityNameSnakeCase: param.EntityNameSnakeCase,
		Location:            param.Location,
		Api:                 param.Api,
	}
}

func (d *usecaseTest) Replace() error {
	paths := []string{
		"templates/usecase_test/usecase_test.tmpl",
	}
	t := template.Must(template.ParseFiles(paths...))

	err := os.MkdirAll(d.Location+"/src/business/usecase/"+d.EntityNameSnakeCase, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(d.Location + "/src/business/usecase/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + "_test.go")
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

func (d *usecaseTest) AppendInterfaceAndFunction() (err error) {
	path := d.Location + "/src/business/usecase/" + d.EntityNameSnakeCase + "/" + d.EntityNameSnakeCase + "_test.go"

	// Read phase: load existing content into memory.
	fRead, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(fRead)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanErr := scanner.Err(); scanErr != nil {
		_ = fRead.Close()
		return scanErr
	}
	if closeErr := fRead.Close(); closeErr != nil {
		return closeErr
	}

	for _, api := range d.Api {
		apiLower := strings.ToLower(api)
		apiLower = strings.ReplaceAll(apiLower, " ", "")
		lines, err = d.appendFunction(lines, apiLower, "templates/usecase_test/usecase_test_temp.tmpl")
		if err != nil {
			return err
		}
	}
	content := strings.Join(lines, "\n")

	// Write phase: re-open the file and rewrite from offset 0. The deferred
	// close preserves any earlier err (CWE-252: writable handles can lose
	// data on Close failure).
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	_, err = f.WriteAt([]byte(content), 0)
	return err
}

func (d *usecaseTest) replaceFunction(templatePath, generatedPath string) (*os.File, error) {
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

func (d *usecaseTest) getFunction(path string) (lines []string, err error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return lines, scanErr
	}
	return lines, nil
}

func (d *usecaseTest) appendFunction(lines []string, api, generatedPath string) ([]string, error) {
	result := lines
	file, err := d.replaceFunction("templates/usecase_test/usecase_test_"+api+"_function.tmpl", generatedPath)
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
