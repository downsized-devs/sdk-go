package parser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/log"
	"github.com/xuri/excelize/v2"
)

const (
	ExcelTypeSow                = "sowexcel"
	ExcelTypePartialHarvest     = "partialharvestexcel"
	ExcelTypeFinalHarvest       = "finalharvestexcel"
	ExcelTypeDailyMonitoring    = "dailymonitoringexcel"
	ExcelTypeSampling           = "samplingexcel"
	ExcelTypeProductionPlan     = "productionplanexcel"
	ExcelTypeFarmTechnician     = "farmtechnicianexcel"
	ExcelTypeLabAnalyst         = "labanalystexcel"
	ExcelTypeOrderBulkAquaCheck = "bulkorderaquacheck"
	ExcelTypeHealthIndex        = "healthindexexcel"
	MetricName                  = "metric"

	ParsingErrMessage      = "err:%+v mapper:%+v row:%d  column: %d "
	ExcelTypeDynamicParser = "dynamic"
	SheetIndexRegex        = `^([a-zA-Z\d]+)(.*)$`
)

var (
	timeLayouts = []string{
		"2-Jan-2006",
		"2-Jan-06",
		"02-Jan-2006",
		"02-Jan-06",
		"01-02-06",
		"2006-01-02",
		"2006-01-2",
		"Jan 2, 06",
		"2 Jan 06",
		"02/01/2006",
		"01/02/2006",
		"01/2/2006",
		"1/02/2006",
		"1/2/2006",
	}
)

type ExcelHeader struct {
	Parent      string
	ParentValue string
	Name        string
	Value       string
	ColIndex    int
}

type ExcelPoint struct {
	Value    string
	RowNum   int
	ColNum   int
	Position string
}

type ExcelResultData map[string]interface{}

type ExcelResultDetails map[string]ExcelPoint

type ExcelResultHeader map[string]ExcelHeader

type ExcelResult struct {
	Header  ExcelResultHeader
	Data    []ExcelResultData
	Details []ExcelResultDetails
}

// TODO: implemented on one template in aquahero, later will be moved here
type ExcelError []struct {
	RowNum   int
	ColNum   int
	Value    string // Old Value
	Position string // "A1, B2, C3"
	Message  string
}

type ExcelOption struct {
	Type      string        `json:"type"`
	Version   string        `json:"version"`
	Parser    Transformer   `json:"parser"`
	Parsers   []Transformer `json:"parsers"`
	ParserMap map[string]Transformer
}

// ExcelInput is used as the structure to generate new excel file
type ExcelInput struct {
	SheetName  string
	SheetValue map[string]string
}

// Note
// For next iteration
// parse column dynamically, read from column name and automatically map it
type ExcelOptions map[string]ExcelOption

type ExcelInterface interface {
	// General Excel
	Unmarshal(ctx context.Context, blob []byte) ([]map[string]string, error)

	// Opiniated Excel
	UnmarshalTransform(ctx context.Context, blob []byte, t Transformer) (ExcelResult, error)
	UnmarshalTransformV2(ctx context.Context, blob []byte, t Transformer) (ExcelResult, error)
	Validate(ctx context.Context, blob []byte, t Transformer) error
	UnmarshalWithAxis(ctx context.Context, blob []byte, t Transformer) (ExcelResult, error)
	UnmarshalTransformMultipleSheet(ctx context.Context, blob []byte, opt ExcelOption) (map[string]ExcelResult, error)
	UnmarshalTransformMultipleSheetV2(ctx context.Context, blob []byte, opt ExcelOption) (map[string]ExcelResult, error)
	UnmarshalTransformYAxis(ctx context.Context, blob []byte, t Transformer) (ExcelResult, error)

	// Generate excel file
	Marshal(ctx context.Context, input []ExcelInput, pathToFile string) error

	WriteReader(ctx context.Context, blob []byte, ioReader io.Writer) error
	InjectValue(ctx context.Context, blob []byte, sheetName string, values map[string]interface{}) ([]byte, error)
}

type excelParser struct {
	log log.Interface
}

func initExcel(log log.Interface) ExcelInterface {
	return &excelParser{
		log: log,
	}
}

// Unmarshal general excel files
func (p *excelParser) Unmarshal(ctx context.Context, blob []byte) ([]map[string]string, error) {
	return nil, errors.NewWithCode(codes.CodeNotImplemented, "Unmarshal not implemented")
}

// Use this to do simple validation, such as header checking, sheet names, etc
func (p *excelParser) Validate(ctx context.Context, blob []byte, t Transformer) error {
	return errors.NewWithCode(codes.CodeNotImplemented, "Validate not implemented")
}

func (p *excelParser) UnmarshalTransform(ctx context.Context, blob []byte, t Transformer) (ExcelResult, error) {
	result := ExcelResult{}
	f, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Slice
	m, err := t.slice(ctx, f)
	p.log.Debug(ctx, len(m))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Filter
	m, err = t.filter(ctx, m)
	p.log.Debug(ctx, len(m))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Map
	result, _, err = t.parse(ctx, m)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	return result, err
}

func (p *excelParser) UnmarshalWithAxis(ctx context.Context, blob []byte, t Transformer) (ExcelResult, error) {
	result := ExcelResult{}
	f, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	sheetExist := f.GetSheetVisible(t.SheetRegex)
	if !sheetExist {
		return result, errors.NewWithCode(codes.CodeExcelFailedParsing, "Invalid sheet name")
	}

	cols, err := f.GetCols(t.SheetRegex)
	if err != nil {
		return result, err
	}

	rows, err := f.GetRows(t.SheetRegex)
	if err != nil {
		return result, err
	}

	if len(cols) != t.SlicerNumCols || len(rows) != t.SlicerNumRows {
		return result, errors.NewWithCode(codes.CodeExcelFailedParsing, "Invalid column or row count")
	}

	for _, v := range t.MappersV2 {
		r := ExcelResultData{}
		rd := ExcelResultDetails{}

		for _, m := range v.Mappers {
			value, err := f.GetCellValue(t.SheetRegex, m.Axis)
			if err != nil {
				return ExcelResult{}, errors.NewWithCode(codes.CodeExcelFailedParsing, fmt.Sprintf("err:%+v mapper:%+v axis:%s", err.Error(), m, m.Axis))
			}

			parsedVal, err := m.ParseString(value)
			if err != nil {
				return ExcelResult{}, errors.NewWithCode(codes.CodeExcelFailedParsing, fmt.Sprintf("err:%+v mapper:%+v axis:%s", err.Error(), m, m.Axis))
			}

			// customize vertical title
			if m.Name == MetricName {
				parsedVal = v.Name
			}

			r[m.Name] = parsedVal
			colNum, rowNum, err := excelize.CellNameToCoordinates(m.Axis)
			if err != nil {
				return result, err
			}

			rd[m.Name] = ExcelPoint{
				Value:    value,
				RowNum:   rowNum,
				ColNum:   colNum,
				Position: m.Axis,
			}
		}

		result.Data = append(result.Data, r)
		result.Details = append(result.Details, rd)
	}

	return result, nil
}

func (p *excelParser) UnmarshalTransformV2(ctx context.Context, blob []byte, transformer Transformer) (ExcelResult, error) {
	result := ExcelResult{}
	excelFile, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Slice
	m, rowLenght, err := transformer.sliceV2(ctx, excelFile)
	p.log.Debug(ctx, len(m))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Map
	result, _, err = transformer.parseSheetDataV2(ctx, rowLenght, m)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	return result, err
}

func (p *excelParser) UnmarshalTransformYAxis(ctx context.Context, blob []byte, transformer Transformer) (ExcelResult, error) {
	result := ExcelResult{}
	excelFile, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Slice
	m, rowLenght, err := transformer.sliceV2(ctx, excelFile)
	p.log.Debug(ctx, len(m))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Map
	result, _, err = transformer.parseRowsDataVertically(ctx, rowLenght, m)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	return result, err
}

func (p *excelParser) UnmarshalTransformMultipleSheet(ctx context.Context, blob []byte, opt ExcelOption) (map[string]ExcelResult, error) {
	result := make(map[string]ExcelResult)
	excelFile, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	opt.initParserMap()

	sheetList := excelFile.GetSheetList()

	// get all list sheet
	for _, sheetName := range sheetList {
		excelSheet := strings.ReplaceAll(strings.ToLower(sheetName), " ", "")
		parser, ok := opt.ParserMap[excelSheet]
		if !ok {
			// Handle dynamic sheets
			if ps, isMatch := opt.getParserForDynamicSheet(sheetName); isMatch {
				parser = ps
			} else {
				continue
			}
		}

		// Parse the sheet data
		var excelResult ExcelResult
		switch parser.SheetType {
		case ExcelTypeDynamicParser:
			var pondCode string
			re := regexp.MustCompile(SheetIndexRegex)
			match := re.FindStringSubmatch(sheetName)
			if len(match) > 1 {
				pondCode = match[1]
			}
			excelResult, _, err = parser.parseSheetDataV1(ctx, excelFile, sheetName)
			if err != nil {
				return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
			}

			// Add the result to the map
			if len(pondCode) > 0 {
				result[fmt.Sprintf("%s-%s", pondCode, parser.Name)] = excelResult
			}
		default:
			// Process the sheet
			matrixResult, headers, rowLenght, err := parser.sliceStaticSheet(ctx, excelFile, sheetName)
			if err != nil {
				return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
			}

			switch {
			case parser.IsDynamicHeader:
				excelResult, _, err = parser.parseDynamicHeadersV2(ctx, rowLenght, matrixResult, headers)
			default:
				excelResult, _, err = parser.parseSheetDataV2(ctx, rowLenght, matrixResult)
			}
			if err != nil {
				return result, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
			}

			// Add the result to the map
			result[parser.Name] = excelResult
		}
	}

	return result, nil
}

func (p *excelParser) UnmarshalTransformMultipleSheetV2(ctx context.Context, blob []byte, opt ExcelOption) (map[string]ExcelResult, error) {
	excelFile, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Init parser map to store parser data
	opt.initParserMap()

	sheetList := excelFile.GetSheetList()

	// Create a channel to send errors
	errChan := make(chan error, len(sheetList))

	// Create a wait group for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(len(sheetList))

	// Create a map to store the results
	results := make(map[string]ExcelResult)

	// Process each sheet concurrently
	for _, sheetName := range sheetList {
		go func(sheetName string) {
			defer wg.Done()

			excelSheet := strings.ReplaceAll(strings.ToLower(sheetName), " ", "")
			parser, ok := opt.ParserMap[excelSheet]
			if !ok {
				// Handle dynamic sheets
				if ps, isMatch := opt.getParserForDynamicSheet(sheetName); isMatch {
					parser = ps
				} else {
					return
				}
			}

			// Parse the sheet data
			var excelResult ExcelResult
			switch parser.SheetType {
			case ExcelTypeDynamicParser:
				var pondCode string
				re := regexp.MustCompile(SheetIndexRegex)
				match := re.FindStringSubmatch(sheetName)
				if len(match) > 1 {
					pondCode = match[1]
				}
				excelResult, _, err = parser.parseSheetDataV1(ctx, excelFile, sheetName)
				if err != nil {
					errChan <- err
					return
				}

				// Add the result to the map
				results[fmt.Sprintf("%s-%s", pondCode, parser.Name)] = excelResult
			default:
				// Process the sheet
				matrixResult, headers, rowLenght, err := parser.sliceStaticSheet(ctx, excelFile, sheetName)
				if err != nil {
					errChan <- err
					return
				}

				switch {
				case parser.IsDynamicHeader:
					excelResult, _, err = parser.parseDynamicHeadersV2(ctx, rowLenght, matrixResult, headers)
				default:
					excelResult, _, err = parser.parseSheetDataV2(ctx, rowLenght, matrixResult)
				}
				if err != nil {
					errChan <- err
					return
				}

				// Add the result to the map
				results[parser.Name] = excelResult
			}
		}(sheetName)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check if there were any errors
	select {
	case err := <-errChan:
		return nil, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	default:
	}

	return results, nil
}

func (p *excelParser) Marshal(ctx context.Context, input []ExcelInput, pathToFile string) error {
	f := excelize.NewFile()
	defer func(ctx context.Context) {
		if err := f.Close(); err != nil {
			p.log.Error(ctx, err)
		}
	}(ctx)

	for i := range input {
		index := f.NewSheet(input[i].SheetName)
		for j, k := range input[i].SheetValue {
			f.SetCellValue(input[i].SheetName, j, k)
		}

		// Set the first sheet as the active sheet
		if i < 1 {
			f.SetActiveSheet(index)
		}
	}

	// Delete default sheet
	f.DeleteSheet("Sheet1")

	// Save spreadsheet to temporary file
	if err := f.SaveAs(pathToFile, excelize.Options{}); err != nil {
		p.log.Error(ctx, err)
		return errors.NewWithCode(codes.CodeExcelFailedToSaveFile, err.Error())
	}

	return nil
}

func (p *excelParser) WriteReader(ctx context.Context, blob []byte, ioReader io.Writer) error {
	excelfile, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return err
	}

	return excelfile.Write(ioReader)
}

func (p *excelParser) InjectValue(ctx context.Context, blob []byte, sheetName string, values map[string]interface{}) ([]byte, error) {
	excelfile, err := excelize.OpenReader(bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	for key, value := range values {
		excelfile.SetCellValue(sheetName, key, value)
	}

	buffer, err := excelfile.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
