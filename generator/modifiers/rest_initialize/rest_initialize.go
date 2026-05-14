package restinItialize

import (
	"bufio"
	"os"
	"strings"
)

type restInitialize struct {
	EntityNameUpper      string
	EntityNameSnakeCase  string
	EntityNameLowerSpace string
	EntityNameLowerDash  string
	Location             string
	Api                  []string
}

type Interface interface {
	Initialize() error
}

type Params struct {
	EntityNameUpper      string
	EntityNameSnakeCase  string
	EntityNameLowerSpace string
	EntityNameLowerDash  string
	Location             string
	Api                  []string
}

func Init(param Params) Interface {
	return &restInitialize{
		EntityNameUpper:      param.EntityNameUpper,
		EntityNameSnakeCase:  param.EntityNameSnakeCase,
		EntityNameLowerSpace: param.EntityNameLowerSpace,
		EntityNameLowerDash:  param.EntityNameLowerDash,
		Location:             param.Location,
		Api:                  param.Api,
	}
}

func (d *restInitialize) Initialize() error {
	f, err := os.OpenFile(d.Location+"/src/handler/rest/rest.go", os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	isFindingRegister := false
	count := 0

	for scanner.Scan() {
		ln := scanner.Text()
		lines, isFindingRegister, count = d.addText(ln, isFindingRegister, count, lines)
	}

	content := strings.Join(lines, "\n")
	_, err = f.WriteAt([]byte(content), 0)
	if err != nil {
		return err
	}
	return nil
}

func (d *restInitialize) addText(ln string, isFindingRegister bool, count int, lines []string) ([]string, bool, int) {
	mapApi := d.mapApi()
	switch {
	case strings.Contains(ln, "func (r *rest) Register() {"):
		isFindingRegister = true
		lines = append(lines, ln)
	case isFindingRegister && strings.Contains(ln, "{"):
		count++
		lines = append(lines, ln)
	case isFindingRegister && strings.Contains(ln, "}") && count == 0:
		temp := []string{"	// " + d.EntityNameLowerSpace}
		if mapApi["create"] {
			temp = append(temp, "	v1.POST(\"/"+d.EntityNameLowerDash+"\", r.Create"+d.EntityNameUpper+")")
		}
		if mapApi["edit"] {
			temp = append(temp, "	v1.PUT(\"/"+d.EntityNameLowerDash+"/:id\", r.Update"+d.EntityNameUpper+")")
		}
		if mapApi["delete"] {
			temp = append(temp, "	v1.DELETE(\"/"+d.EntityNameLowerDash+"/:id\", r.AuthorizeScope(entity.AdminDelete, r.Delete"+d.EntityNameUpper+"))")
		}
		if mapApi["get"] {
			temp = append(temp, "	v1.GET(\"/admin/"+d.EntityNameLowerDash+"\", r.AuthorizeScope(entity.AdminList, r.Get"+d.EntityNameUpper+"ListForAdmin))")
		}
		if mapApi["activate"] {
			temp = append(temp, "	v1.PATCH(\"/admin/"+d.EntityNameLowerDash+"/:id/activate\", r.AuthorizeScope(entity.AdminActivate, r.Activate"+d.EntityNameUpper+"))")
		}

		lines = append(lines, temp...)
		lines = append(lines, ln)
	case isFindingRegister && strings.Contains(ln, "}"):
		count--
		lines = append(lines, ln)
	default:
		lines = append(lines, ln)
	}

	return lines, isFindingRegister, count
}

func (d *restInitialize) mapApi() map[string]bool {
	result := map[string]bool{
		"create":   false,
		"edit":     false,
		"get":      false,
		"delete":   false,
		"activate": false,
	}
	for _, v := range d.Api {
		apiLower := strings.ToLower(v)
		apiLower = strings.ReplaceAll(apiLower, " ", "")
		result[apiLower] = true
	}
	return result
}
