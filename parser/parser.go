package parser

import (
	"github.com/downsized-devs/sdk-go/log"
)

type Parser interface {
	JSONParser() JSONInterface
	CSVParser() CSVInterface
	ExcelParser() ExcelInterface
}

type Options struct {
	JSONOptions  JSONOptions
	CSVOptions   CSVOptions
	ExcelOptions ExcelOptions
}

type parser struct {
	json  JSONInterface
	csv   CSVInterface
	excel ExcelInterface
}

func InitParser(log log.Interface, opt Options) Parser {
	return &parser{
		json:  initJSON(opt.JSONOptions, log),
		csv:   initCSV(),
		excel: initExcel(log),
	}
}

func (p *parser) JSONParser() JSONInterface {
	return p.json
}

func (p *parser) CSVParser() CSVInterface {
	return p.csv
}

func (p *parser) ExcelParser() ExcelInterface {
	return p.excel
}
