package parser

import (
	"bytes"
	"encoding/csv"
	"io"

	"github.com/gocarina/gocsv"
)

type CsvOptions struct {
	Separator  rune
	LazyQuotes bool
}

type CsvInterface interface {
	SetOptionsRead(opt CsvOptions)
	SetOptionsWrite(opt CsvOptions)
	Marshal(orig interface{}) ([]byte, error)
	MarshalWithoutHeaders(orig interface{}) ([]byte, error)
	Unmarshal(blob []byte, dest interface{}) error
}

type csvParser struct{}

func initCsv() CsvInterface {
	return &csvParser{}
}

func (p *csvParser) Unmarshal(blob []byte, dest interface{}) error {
	err := gocsv.UnmarshalBytes(blob, dest)
	return err
}

func (p *csvParser) Marshal(orig interface{}) (result []byte, err error) {
	result, err = gocsv.MarshalBytes(orig)
	return result, err
}

func (p *csvParser) MarshalWithoutHeaders(orig interface{}) (result []byte, err error) {
	buf := bytes.NewBuffer(nil)
	err = gocsv.MarshalCSVWithoutHeaders(orig, gocsv.DefaultCSVWriter(buf))
	return buf.Bytes(), err
}

func (p *csvParser) SetOptionsRead(opt CsvOptions) {
	gocsv.TagSeparator = string(opt.Separator)
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = opt.LazyQuotes
		r.Comma = opt.Separator
		return r
	})
}

func (p *csvParser) SetOptionsWrite(opt CsvOptions) {
	gocsv.TagSeparator = string(opt.Separator)
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		w := csv.NewWriter(out)
		w.Comma = opt.Separator
		return gocsv.NewSafeCSVWriter(w)
	})
}
