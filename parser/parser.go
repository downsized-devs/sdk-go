package parser

import "github.com/downsized-devs/sdk-go/logger"

type Parser interface {
	JsonParser() JsonInterface
	CsvParser() CsvInterface
	ExcelParser() ExcelInterface
}

type Options struct {
	JsonOptions  JsonOptions
	CsvOptions   CsvOptions
	ExcelOptions ExcelOptions
}

type parser struct {
	Json  JsonInterface
	Csv   CsvInterface
	excel ExcelInterface
}

func InitParser(log logger.Interface, opt Options) Parser {
	return &parser{
		Json:  initJson(opt.JsonOptions, log),
		Csv:   initCsv(),
		excel: initExcel(log),
	}
}

func (p *parser) JsonParser() JsonInterface {
	return p.Json
}

func (p *parser) CsvParser() CsvInterface {
	return p.Csv
}

func (p *parser) ExcelParser() ExcelInterface {
	return p.excel
}
