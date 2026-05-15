package parser

import "github.com/downsized-devs/sdk-go/logger"

type Interface interface {
	JsonParser() JsonInterface
	CsvParser() CsvInterface
}

type Options struct {
	JsonOptions JsonOptions
	CsvOptions  CsvOptions
}

type parser struct {
	Json JsonInterface
	Csv  CsvInterface
}

func Init(log logger.Interface, opt Options) Interface {
	return &parser{
		Json: initJson(opt.JsonOptions, log),
		Csv:  initCsv(),
	}
}

func (p *parser) JsonParser() JsonInterface {
	return p.Json
}

func (p *parser) CsvParser() CsvInterface {
	return p.Csv
}
