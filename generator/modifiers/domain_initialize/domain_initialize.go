package domaininItialize

import (
	"bufio"
	"os"
	"strings"
)

type domainInitialize struct {
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
	return &domainInitialize{
		EntityNameUpper:     param.EntityNameUpper,
		EntityNameSnakeCase: param.EntityNameSnakeCase,
		Location:            param.Location,
	}
}

func (d *domainInitialize) Initialize() error {
	f, err := os.OpenFile(d.Location+"/src/business/domain/domain.go", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		ln := scanner.Text()
		lines = append(lines, ln)
		switch {
		case strings.Contains(ln, "import ("):
			lines = append(lines, "	\"github.com/downsized-devs/generic-service/src/business/domain/"+d.EntityNameSnakeCase+"\"")
			continue
		case strings.Contains(ln, "type Domains struct {"):
			lines = append(lines, "	"+d.EntityNameUpper+" "+d.EntityNameSnakeCase+".Interface")
		case strings.Contains(ln, "d := &Domains{"):
			lines = append(lines, "		"+d.EntityNameUpper+":           "+d.EntityNameSnakeCase+".Init(params.Log, params.Db, params.Rd, params.Json, params.Audit),")
		}
	}

	content := strings.Join(lines, "\n")
	_, err = f.WriteAt([]byte(content), 0)
	if err != nil {
		return err
	}
	return nil
}
