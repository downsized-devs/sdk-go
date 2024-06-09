package usecaseInitialize

import (
	"bufio"
	"os"
	"strings"
)

type usecaseInitialize struct {
	EntityNameUpper     string
	EntityNameSnakeCase string
	Location            string
}

type Interface interface {
	Initialize() error
}

type Params struct {
	EntityNameUpper     string
	EntityNameSnakeCase string
	Location            string
}

func Init(param Params) Interface {
	return &usecaseInitialize{
		EntityNameUpper:     param.EntityNameUpper,
		EntityNameSnakeCase: param.EntityNameSnakeCase,
		Location:            param.Location,
	}
}

func (d *usecaseInitialize) Initialize() error {
	f, err := os.OpenFile(d.Location+"/src/business/usecase/usecase.go", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		ln := scanner.Text()
		lines = append(lines, ln)
		if strings.Contains(ln, "import (") {
			lines = append(lines, "	\"github.com/downsized-devs/generic-service/src/business/usecase/"+d.EntityNameSnakeCase+"\"")
			continue
		} else if strings.Contains(ln, "type Usecases struct {") {
			lines = append(lines, "	"+d.EntityNameUpper+" "+d.EntityNameSnakeCase+".Interface")
		} else if strings.Contains(ln, "uc := &Usecases{") {
			lines = append(lines, "		"+d.EntityNameUpper+":           "+d.EntityNameSnakeCase+".Init("+d.EntityNameSnakeCase+".InitParam{Log: params.Log, Dom: params.Dom."+d.EntityNameUpper+", Auth: params.Auth}),")
		}
	}

	content := strings.Join(lines, "\n")
	_, err = f.WriteAt([]byte(content), 0)
	if err != nil {
		return err
	}
	return nil
}
