package email

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Boostport/mjml-go"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
)

type TemplateInterface interface {
	FromHTML(params BodyFromHTMLParams) (string, error)
	FromMJML(ctx context.Context, params BodyFromMJMLParams) (string, error)
}

type emailtemplate struct {
	config   TemplateConfig
	log      logger.Interface
	template map[string]*template.Template
}

type TemplateConfig struct {
	FileDirectory string
}

func initTemplate(conf TemplateConfig, log logger.Interface) TemplateInterface {
	return &emailtemplate{
		config:   conf,
		log:      log,
		template: make(map[string]*template.Template),
	}
}

func (t *emailtemplate) FromHTML(params BodyFromHTMLParams) (string, error) {
	filePath := params.Filename
	if t.config.FileDirectory != "" {
		filePath = fmt.Sprintf("%s/%s", t.config.FileDirectory, filePath)
	}
	if params.OverrideDirectory != "" {
		filePath = fmt.Sprintf("%s/%s", params.OverrideDirectory, params.Filename)
	}

	var err error
	tmpl, ok := t.template[filePath]
	if !ok {
		tmpl, err = template.New(params.Filename).Funcs(template.FuncMap(params.FuncMap)).ParseFiles(filePath)
		if err != nil {
			return "", errors.NewWithCode(codes.CodeParseHTMlTemplateFailed, err.Error())
		}

		t.template[filePath] = tmpl
	}

	body := new(bytes.Buffer)
	if err := tmpl.Execute(body, params.Data); err != nil {
		return "", errors.NewWithCode(codes.CodeParseHTMlTemplateFailed, err.Error())
	}

	return body.String(), nil
}

func (t *emailtemplate) FromMJML(ctx context.Context, params BodyFromMJMLParams) (string, error) {
	filePath := params.Filename
	if t.config.FileDirectory != "" {
		filePath = fmt.Sprintf("%s/%s", t.config.FileDirectory, filePath)
	}
	if params.OverrideDirectory != "" {
		filePath = fmt.Sprintf("%s/%s", params.OverrideDirectory, params.Filename)
	}

	var err error
	tmpl, ok := t.template[filePath]
	if !ok {
		tmpl, err = template.New(params.Filename).Funcs(template.FuncMap(params.FuncMap)).ParseFiles(filePath)
		if err != nil {
			return "", errors.NewWithCode(codes.CodeParseHTMlTemplateFailed, err.Error())
		}

		t.template[filePath] = tmpl
	}

	bodyMjml := new(bytes.Buffer)
	if err := tmpl.Execute(bodyMjml, params.Data); err != nil {
		return "", errors.NewWithCode(codes.CodeConvertMJMLToHTMLFailed, err.Error())
	}

	var mjmlErr mjml.Error
	body, err := mjml.ToHTML(ctx, bodyMjml.String(), mjml.WithMinify(true))
	if err != nil {
		if errors.As(err, &mjmlErr) {
			t.log.Error(ctx, mjmlErr.Message)
			t.log.Error(ctx, mjmlErr.Details)
		}
		return "", errors.NewWithCode(codes.CodeConvertMJMLToHTMLFailed, err.Error())
	}

	return body, nil
}
