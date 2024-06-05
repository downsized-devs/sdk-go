package parser

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cstockton/go-conv"
	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/convert"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/xuri/excelize/v2"
)

type Transformer struct {
	Name            string
	SheetType       string
	SheetName       string
	SheetRegex      string
	SheetSeparators string

	SlicerInitialSkipRows int
	SlicerInitialSkipCols int
	SlicerNumRows         int
	SlicerNumCols         int
	IsDynamicHeader       bool
	Filters               []Filter

	Mappers []Mapper // Beware it's not exactly same number with the excel col number

	MappersV2   map[string]MapperV2
	MapperYAxis MapperYAxis
}

type MapperYAxis struct {
	KeysColumn []string
	MappersVal map[string]Mapper
}

type Filter struct {
	ColNum int
	Regex  string // regex -> pass data when regex is matched
}

type Mapper struct {
	Name        string
	ColNum      int
	Type        string
	Relation    string
	Axis        string
	RelationMap map[string]int64
}

type MapperV2 struct {
	Name        string
	Type        string
	Relation    string
	RelationMap map[string]int64
	Multiple    bool
	MultiColumn bool
	FilterRegex string
	Mappers     []Mapper
}

func (t *Transformer) slice(ctx context.Context, f *excelize.File) ([][]ExcelPoint, error) {
	result := [][]ExcelPoint{}
	var err error

	for _, s := range f.GetSheetList() {
		match, err := regexp.MatchString(t.SheetRegex, s)
		if err != nil {
			return nil, err
		}
		if match {
			if t.SheetSeparators != "" {
				sheets := strings.Split(s, t.SheetSeparators)
				if len(sheets) < 1 {

					return nil, errors.NewWithCode(codes.CodeExcelFailedParsing, "sheet separators Failed")
				}

				for _, sheetSingle := range sheets {
					r, err := t.getRowsAndSlice(ctx, f, s, []string{sheetSingle})
					if err != nil {
						return nil, err
					}
					result = append(result, r...)
				}

			} else {
				appender := []string{}
				if t.SheetType == ExcelTypeDynamicParser {
					appender = []string{s}
				}
				r, err := t.getRowsAndSlice(ctx, f, s, appender)
				if err != nil {
					return nil, err
				}
				result = append(result, r...)
			}
		}
	}

	return result, err

}

func (t *Transformer) filter(ctx context.Context, m [][]ExcelPoint) ([][]ExcelPoint, error) {
	if len(t.Filters) < 1 {
		return m, nil
	}

	result := [][]ExcelPoint{}
	for _, f := range t.Filters {
		for _, r := range m {
			match, err := regexp.MatchString(f.Regex, r[f.ColNum].Value)

			if err != nil {
				return nil, err
			}

			if match {
				result = append(result, r)
			}

		}
	}

	return result, nil

}

func (t *Transformer) parse(ctx context.Context, m [][]ExcelPoint) (ExcelResult, []string, error) {
	var parsingErr []string
	var result ExcelResult

	for _, row := range m {
		r := ExcelResultData{}
		rd := ExcelResultDetails{}
		for _, m := range t.Mappers {

			// This is unsafe :\
			if row[m.ColNum].Value == "" {
				continue
			}
			parsedValue, err := m.ParseString(row[m.ColNum].Value)
			if err != nil {
				parsingErr = append(parsingErr, fmt.Sprintf(ParsingErrMessage, err.Error(), m, row[m.ColNum].RowNum, row[m.ColNum].ColNum))
				continue
			}
			r[m.Name] = parsedValue
			row[m.ColNum].Position = fmt.Sprintf("%s%d", convert.IntToChar(row[m.ColNum].ColNum+1), row[m.ColNum].RowNum+1)
			rd[m.Name] = row[m.ColNum]

		}
		result.Data = append(result.Data, r)
		result.Details = append(result.Details, rd)

	}

	return result, parsingErr, nil
}

func (t *Transformer) getRowsAndSlice(ctx context.Context, f *excelize.File, sheetName string, colAdders []string) ([][]ExcelPoint, error) {
	result := [][]ExcelPoint{}
	fixMergedRow(ctx, f, sheetName)

	rows, err := parseRows(ctx, f, sheetName)

	if err != nil {
		return nil, err
	}

	// check if row count match slicer config to avoid panic
	if len(rows) < t.SlicerInitialSkipRows+1 {
		return nil, errors.NewWithCode(codes.CodeExcelFailedParsing, "invalid row count")
	}

	// Cut x first rows
	rows = rows[t.SlicerInitialSkipRows+1:]

	// rows = rows[:10] // CUTTER PETER
	rows = filterEmptyRows(ctx, rows)
	rows = cutMatrix(ctx, rows, t.SlicerNumCols)
	adders := []ExcelPoint{}
	for _, a := range colAdders {
		adders = append(adders, ExcelPoint{
			Value: a,
		})
	}

	for i, r := range rows {
		r = append(r, adders...)
		rows[i] = r
	}

	result = append(result, rows...)
	return result, nil
}

func (t *Transformer) sliceStaticSheet(ctx context.Context, excelFile *excelize.File, sheetName string) (map[string][]ExcelPoint, []ExcelHeader, int, error) {
	if !regexp.MustCompile(t.SheetRegex).MatchString(sheetName) {
		return nil, nil, 0, nil
	}

	var appender []string
	if t.SheetType == ExcelTypeDynamicParser {
		appender = []string{sheetName}
	}

	return t.getRowsAndSliceV2(ctx, excelFile, sheetName, appender)
}

// Slicing data into map[string][]ExcelPoint, string will be header
// i.e: map{"pondaddress":["A1","A2"]}
func (t *Transformer) sliceV2(ctx context.Context, excelFile *excelize.File) (map[string][]ExcelPoint, int, error) {
	result := map[string][]ExcelPoint{}
	var err error
	var rowLenght = 0

	for _, excelSheet := range excelFile.GetSheetList() {
		match, err := regexp.MatchString(t.SheetRegex, excelSheet)
		if err != nil {
			return nil, 0, err
		}
		if match {
			if t.SheetSeparators != "" {
				sheets := strings.Split(excelSheet, t.SheetSeparators)
				if len(sheets) < 1 {
					return nil, 0, errors.NewWithCode(codes.CodeExcelFailedParsing, "sheet separators Failed")
				}

				for _, sheetSingle := range sheets {
					result, _, rowLenght, err = t.getRowsAndSliceV2(ctx, excelFile, excelSheet, []string{sheetSingle})
					if err != nil {
						return nil, 0, err
					}
				}
			} else {
				appender := []string{}
				if t.SheetType == ExcelTypeDynamicParser {
					appender = []string{excelSheet}
				}
				result, _, rowLenght, err = t.getRowsAndSliceV2(ctx, excelFile, excelSheet, appender)
				if err != nil {
					return nil, 0, err
				}
			}
		}
	}
	return result, rowLenght, err
}

func (t *Transformer) parseSheetDataV1(ctx context.Context, excelFile *excelize.File, sheetName string) (ExcelResult, int, error) {
	// Slice
	rows, err := t.getRowsAndSlice(ctx, excelFile, sheetName, []string{})
	if err != nil {
		return ExcelResult{}, 0, err
	}

	// Filter
	rows, err = t.filter(ctx, rows)
	if err != nil {
		return ExcelResult{}, 0, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	// Map
	parseResult, _, err := t.parse(ctx, rows)
	if err != nil {
		return ExcelResult{}, 0, errors.NewWithCode(codes.CodeClientErrorOnReadBody, err.Error())
	}

	return parseResult, len(rows), nil
}

// getRowsAndSliceV2 will do main slicing operation
func (t *Transformer) getRowsAndSliceV2(ctx context.Context, excelFile *excelize.File, sheetName string, colAdders []string) (map[string][]ExcelPoint, []ExcelHeader, int, error) {
	fixMergedRow(ctx, excelFile, sheetName)

	rows, err := parseRows(ctx, excelFile, sheetName)
	if err != nil {
		return nil, []ExcelHeader{}, 0, err
	}

	if t.IsDynamicHeader {
		return t.getDynamicHeaderResult(ctx, rows)
	}

	return t.getStaticHeaderResult(ctx, rows)
}

func (t *Transformer) getDynamicHeaderResult(ctx context.Context, rows [][]ExcelPoint) (map[string][]ExcelPoint, []ExcelHeader, int, error) {
	result := map[string][]ExcelPoint{}
	// Get headers to fill in map data
	headers, lastCol := t.getDynamicHeaders(ctx, rows)
	// check if row count match slicer config to avoid panic
	if len(rows) < t.SlicerInitialSkipRows+2 {
		return nil, headers, 0, errors.NewWithCode(codes.CodeExcelFailedParsing, "invalid row count")
	}

	// Cut x first rows
	rows = rows[t.SlicerInitialSkipRows+2:]

	rows = filterEmptyRows(ctx, rows)
	rows = cutMatrix(ctx, rows, lastCol)

	for _, r := range rows {
		for _, header := range headers {
			if header.Parent != header.Name {
				key := fmt.Sprintf("%s:%s", header.Parent, header.Name)
				result[key] = append(result[key], r[header.ColIndex])
			} else {
				result[header.Name] = append(result[header.Name], r[header.ColIndex])
			}

		}
	}

	return result, headers, len(rows), nil
}

func (t *Transformer) getStaticHeaderResult(ctx context.Context, rows [][]ExcelPoint) (map[string][]ExcelPoint, []ExcelHeader, int, error) {
	result := map[string][]ExcelPoint{}
	// Get headers to fill in map data
	headers, lastCol := t.getSingleHeaders(ctx, rows)
	// check if row count match slicer config to avoid panic
	if len(rows) < t.SlicerInitialSkipRows+1 {
		return nil, []ExcelHeader{}, 0, errors.NewWithCode(codes.CodeExcelFailedParsing, "invalid row count")
	}

	// Cut x first rows
	rows = rows[t.SlicerInitialSkipRows+1:]

	rows = filterEmptyRows(ctx, rows)
	rows = cutMatrix(ctx, rows, lastCol)

	for _, r := range rows {
		for _, header := range headers {
			result[header.Name] = append(result[header.Name], r[header.ColIndex])
		}
	}

	return result, headers, len(rows), nil
}

func (t *Transformer) getDynamicHeaders(ctx context.Context, rows [][]ExcelPoint) ([]ExcelHeader, int) {
	headers := []ExcelHeader{}
	lastCol := 0

	newRows := rows[t.SlicerInitialSkipRows : t.SlicerInitialSkipRows+2]
	// newRows[0] ==> is parent header
	// newRows[1] ==> is child header

	for i, row := range newRows[1] {
		header := ExcelHeader{
			// name will be made lowercase and all space will be removed for more flexibility
			Parent:      strings.ReplaceAll(strings.ToLower(newRows[0][i].Value), " ", ""),
			ParentValue: newRows[0][i].Value,
			Name:        strings.ReplaceAll(strings.ToLower(row.Value), " ", ""),
			Value:       row.Value,
			ColIndex:    row.ColNum,
		}

		headers = append(headers, header)
	}

	lastCol = len(newRows[1])
	return headers, lastCol
}

func (t *Transformer) getSingleHeaders(ctx context.Context, rows [][]ExcelPoint) ([]ExcelHeader, int) {
	headers := []ExcelHeader{}
	lastCol := 0

	newRows := rows[t.SlicerInitialSkipRows : t.SlicerInitialSkipRows+1]
	for _, row := range newRows[0] {
		header := ExcelHeader{
			// name will be made lowercase and all space will be removed for more flexibility
			Name:     strings.ReplaceAll(strings.ToLower(row.Value), " ", ""),
			Value:    row.Value,
			ColIndex: row.ColNum,
		}

		headers = append(headers, header)
	}

	lastCol = len(newRows[0])

	return headers, lastCol
}

func (t *Transformer) filterExcelData(ctx context.Context, r ExcelResultData) bool {
	for _, m := range t.MappersV2 {
		if m.FilterRegex != "" {
			//string conversion
			v, err := convert.ToString(r[m.Name])
			if err != nil {
				return false
			}

			match, err := regexp.MatchString(m.FilterRegex, v)
			if err != nil || !match {
				return false
			}
		}
	}

	return true
}

func (m *Mapper) ParseString(input string) (interface{}, error) {
	switch m.Type {
	case "string":
		if m.Relation != "" {
			id, ok := m.RelationMap[input]
			if !ok {
				return 0, errors.NewWithCode(codes.CodeExcelFailedParsing, "failed making relation")
			}
			return id, nil
		}
		return input, nil
	case "int":
		input = strings.TrimSpace(input)
		r, err := conv.Int64(input)
		return r, err
	case "float":
		input = strings.TrimSpace(input)
		r, err := conv.Float64(input)
		return r, err
	case "dynamic_date":
		var err error
		input = strings.TrimSpace(input)
		for _, layout := range timeLayouts {
			t, err := time.Parse(layout, input)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, err
	case "date:02-Jan-2006":
		layout := "2-Jan-2006"
		t, err := time.Parse(layout, input)
		if err != nil {
			return t, err
		}
		return t, nil
	case "date:02-Jan-06":
		layout := "2-Jan-06"
		t, err := time.Parse(layout, input)
		if err != nil { // FU
			layout := "2-Jan-2006"
			t, err = time.Parse(layout, input)
			if err != nil {
				return t, err
			}
			return t, nil
		}
		return t, nil
	case "date:01-02-06":
		layout := "01-02-06"
		t, err := time.Parse(layout, input)
		if err != nil {
			return t, err
		}
		return t, nil
	case "date:2006-01-02":
		layout := "2006-01-02"
		t, err := time.Parse(layout, input)
		if err != nil {
			return t, err
		}
		return t, nil
	default:
		return input, nil
	}
}

func (m MapperV2) ParseStruct(name string, input []ExcelPoint) (interface{}, error) {

	res := make(map[string]interface{})

	for _, mapper := range m.Mappers {
		parsedValue, err := mapper.ParseString(input[mapper.ColNum].Value)
		if err != nil {
			continue
		}

		res[mapper.Name] = parsedValue
	}

	if name != "" {
		res["name"] = name
	}

	return res, nil
}

func (m MapperV2) ParseString(input string) (interface{}, error) {
	switch m.Type {
	case "string":
		if m.Relation != "" {
			id, ok := m.RelationMap[input]
			if !ok {
				return 0, errors.NewWithCode(codes.CodeExcelFailedParsing, "failed making relation")
			}
			return id, nil
		}
		return input, nil
	case "int":
		r, err := conv.Int64(input)
		return r, err
	case "float":
		r, err := conv.Float64(input)
		return r, err
	case "date":
		//TODO
		return input, nil

	case "dynamic_date":
		var err error
		for _, layout := range timeLayouts {
			t, err := time.Parse(layout, input)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, err
	case "date:02-Jan-2006":
		layout := "2-Jan-2006"
		t, err := time.Parse(layout, input)
		if err != nil {
			return t, err
		}
		return t, nil
	case "date:02-Jan-06":
		layout := "2-Jan-06"
		t, err := time.Parse(layout, input)
		if err != nil { // FU
			layout := "2-Jan-2006"
			t, err = time.Parse(layout, input)
			if err != nil {
				return t, err
			}
			return t, nil
		}
		return t, nil
	case "date:01-02-06":
		layout := "01-02-06"
		t, err := time.Parse(layout, input)
		if err != nil {
			return t, err
		}
		return t, nil
	case "date:2006-01-02":
		layout := "2006-01-02"
		t, err := time.Parse(layout, input)
		if err != nil {
			return t, err
		}
		return t, nil
	default:
		return input, nil
	}
}

// TODO Unit test
func filterEmptyRows(ctx context.Context, m [][]ExcelPoint) [][]ExcelPoint {
	result := [][]ExcelPoint{}
	for _, row := range m {
		if len(row) >= 1 {
			result = append(result, row)
		}
	}
	return result
}

// TODO Unit test
func cutMatrix(ctx context.Context, m [][]ExcelPoint, nCol int) [][]ExcelPoint {
	result := [][]ExcelPoint{}

	for _, row := range m {
		l := len(row)
		if l > nCol {
			result = append(result, row[:nCol])
		}
		if l < nCol {
			row = append(row, EmptyRow(nCol-l)...)
			result = append(result, row)
		} else {
			result = append(result, row)
		}

	}

	return result
}

func EmptyRow(nCol int) []ExcelPoint {
	result := []ExcelPoint{}
	for i := 0; i < nCol; i++ {
		result = append(result, ExcelPoint{})
	}
	return result
}

// TODO: show what exactly makes the file error
// func showErrorInExcel(ctx context.Context, f *excelize.File) {}

func fixMergedRow(ctx context.Context, f *excelize.File, sheetName string) error {
	mergeCells, err := f.GetMergeCells(sheetName)
	if err != nil {
		return err
	}

	for _, m := range mergeCells {
		coord := strings.Split(m[0], ":")
		start := coord[0]
		end := coord[1]
		f.UnmergeCell(sheetName, start, end)
	}

	for _, m := range mergeCells {
		coordinates, err := generateCoords(ctx, m[0])
		if err != nil {
			return err
		}
		for _, coordinate := range coordinates {
			f.SetCellValue(sheetName, coordinate, m[1])
		}
	}
	return nil
}

// i.e. D10:D13 -> [D10, D11, D12, D13]
func generateCoords(ctx context.Context, coordinateRange string) ([]string, error) {
	result := []string{}
	coords := strings.Split(coordinateRange, ":")
	startCol, startRow, err := excelize.CellNameToCoordinates(coords[0])
	if err != nil {
		return result, err
	}
	endCol, endRow, err := excelize.CellNameToCoordinates(coords[1])
	if err != nil {
		return result, err
	}

	for c := startCol; c <= endCol; c++ {
		for r := startRow; r <= endRow; r++ {
			namedCoor, err := excelize.CoordinatesToCellName(c, r)
			if err != nil {
				return result, err
			}
			result = append(result, namedCoor)
		}
	}

	return result, nil
}

func parseRows(ctx context.Context, f *excelize.File, sheetName string) ([][]ExcelPoint, error) {
	result := [][]ExcelPoint{}
	rows, err := f.GetRows(sheetName)

	if err != nil {
		return nil, err
	}

	for rowNum, row := range rows {
		rowResult := []ExcelPoint{}
		for colNum, cell := range row {
			p := ExcelPoint{
				Value:  cell,
				RowNum: rowNum,
				ColNum: colNum,
			}
			rowResult = append(rowResult, p)

		}
		result = append(result, rowResult)
	}
	return result, nil
}

func (opt *ExcelOption) initParserMap() {
	if opt.ParserMap == nil {
		opt.ParserMap = make(map[string]Transformer)
	}

	for _, parser := range opt.Parsers {
		opt.ParserMap[strings.ReplaceAll(strings.ToLower(parser.SheetRegex), " ", "")] = parser
	}
}

func (opt *ExcelOption) getParserForDynamicSheet(sheetName string) (Transformer, bool) {
	for _, parser := range opt.Parsers {
		if parser.SheetType == ExcelTypeDynamicParser && regexp.MustCompile(parser.SheetRegex).MatchString(sheetName) {
			return parser, true
		}
	}

	return Transformer{}, false
}

func (t *Transformer) parseDynamicHeadersV2(ctx context.Context, rowLength int, rows map[string][]ExcelPoint, headers []ExcelHeader) (ExcelResult, []string, error) {
	var parsingErr []string
	result := ExcelResult{Header: ExcelResultHeader{}}

	for _, header := range headers {
		if mapper, ok := t.MappersV2[header.Name]; ok {
			result.Header[mapper.Name] = header
		} else {
			result.Header[header.Parent+":"+header.Name] = header
		}
	}

	for i := 0; i < rowLength; i++ {
		rowData := ExcelResultData{}
		rowDetails := ExcelResultDetails{}
		for columnName, excelPoints := range rows {
			mapper, ok := t.MappersV2[columnName]
			if !ok {
				t.MappersV2[columnName] = MapperV2{
					Name: columnName,
					Type: reflect.TypeOf(excelPoints[i].Value).String(),
				}
				mapper = t.MappersV2[columnName]
			}

			parsedValue, err := mapper.ParseString(excelPoints[i].Value)
			if err != nil {
				parsingErr = append(parsingErr, fmt.Sprintf(ParsingErrMessage, err.Error(), mapper.Name, excelPoints[i].RowNum, excelPoints[i].ColNum))
			}

			rowData[mapper.Name] = parsedValue
			rowDetails[mapper.Name] = excelPoints[i]
		}

		if t.filterExcelData(ctx, rowData) {
			result.Data = append(result.Data, rowData)
			result.Details = append(result.Details, rowDetails)
		}
	}

	return result, parsingErr, nil
}

func (t *Transformer) parseSheetDataV2(ctx context.Context, rowLength int, rows map[string][]ExcelPoint) (ExcelResult, []string, error) {
	var parsingErrs []string
	var result ExcelResult

	for i := 0; i < rowLength; i++ {
		rowData, rowDetails, errs := t.parseRows(rows, i)
		if len(errs) > 0 {
			parsingErrs = append(parsingErrs, errs...)
		}

		if t.filterExcelData(ctx, rowData) {
			result.Data = append(result.Data, rowData)
			result.Details = append(result.Details, rowDetails)
		}
	}

	return result, parsingErrs, nil
}

func (t *Transformer) parseRowsDataVertically(ctx context.Context, rowLength int, rows map[string][]ExcelPoint) (ExcelResult, []string, error) {
	var result ExcelResult
	parsingErr := []string{}
	mapFinExcel := make(map[string]map[string]interface{})
	keysColumn := []string{}

	for indexRow := 0; indexRow < rowLength; indexRow++ {
		key := ""
		for j, columnName := range t.MapperYAxis.KeysColumn {
			if j > 0 {
				key += "-"
			}
			key += rows[columnName][indexRow].Value
		}

		if mapper, isExist := t.MapperYAxis.MappersVal[key]; isExist {
			for columnName, excelPoints := range rows {
				isSkip := t.isIncludeKeyColumn(columnName)
				if isSkip {
					continue
				}

				_, isMatched := t.isKeyExistWithRegex(columnName)
				if !isMatched {
					continue
				}

				parsedVal, err := mapper.ParseString(excelPoints[indexRow].Value)
				if err != nil {
					parsingErr = append(parsingErr, fmt.Sprintf(ParsingErrMessage, err.Error(), columnName, excelPoints[indexRow].RowNum, excelPoints[indexRow].ColNum))
				}

				if _, isColumnExist := mapFinExcel[columnName]; isColumnExist {
					mapFinExcel[columnName][mapper.Name] = parsedVal
				} else {
					mapFinExcel[columnName] = map[string]interface{}{
						mapper.Name: parsedVal,
					}
				}
			}
		}
	}

	for column, _ := range mapFinExcel {
		keysColumn = append(keysColumn, column)
	}

	sort.Strings(keysColumn)
	for _, key := range keysColumn {
		value := mapFinExcel[key]
		value["keyColumn"] = key
		result.Data = append(result.Data, value)
	}

	return result, parsingErr, nil
}

func (t *Transformer) isIncludeKeyColumn(key string) bool {
	for _, s := range t.MapperYAxis.KeysColumn {
		if s == key {
			return true
		}
	}
	return false
}

func (t *Transformer) parseRows(rows map[string][]ExcelPoint, rowIndex int) (ExcelResultData, ExcelResultDetails, []string) {
	rowData := ExcelResultData{}
	rowDetails := ExcelResultDetails{}
	parsingErr := []string{}

	for columnName, excelPoints := range rows {
		mapper, isMatched := t.isKeyExistWithRegex(columnName)
		if !isMatched {
			mapper = t.MappersV2[columnName]
		}

		parsedValue, err := t.parseColumn(excelPoints, rowIndex, mapper, columnName)
		if err != nil {
			parsingErr = append(parsingErr, fmt.Sprintf(ParsingErrMessage, err.Error(), columnName, excelPoints[rowIndex].RowNum, excelPoints[rowIndex].ColNum))
		}

		if mapper.Multiple {
			if list, ok := rowData[mapper.Name].([]interface{}); ok {
				rowData[mapper.Name] = append(list, parsedValue)
			} else {
				rowData[mapper.Name] = []interface{}{parsedValue}
			}
		} else {
			rowData[mapper.Name] = parsedValue
			rowDetails[mapper.Name] = excelPoints[rowIndex]
		}
	}

	return rowData, rowDetails, parsingErr
}

func (t *Transformer) isKeyExistWithRegex(columnName string) (MapperV2, bool) {
	for mapperKey, mapperValue := range t.MappersV2 {
		if mapperValue.Multiple {
			if matched, _ := regexp.MatchString(mapperKey, columnName); matched {
				return mapperValue, true
			}
		}
	}

	return MapperV2{}, false
}

func (t *Transformer) parseColumn(excelPoints []ExcelPoint, rowIndex int, mapper MapperV2, columnName string) (interface{}, error) {
	if mapper.MultiColumn {
		return t.parseMultiColumn(excelPoints, rowIndex, mapper, columnName)
	} else {
		return t.parseSingleColumn(excelPoints, rowIndex, mapper, columnName)
	}
}

func (t *Transformer) parseSingleColumn(excelPoints []ExcelPoint, rowIndex int, mapper MapperV2, columnName string) (interface{}, error) {
	return mapper.ParseString(excelPoints[rowIndex].Value)
}

func (t *Transformer) parseMultiColumn(excelPoints []ExcelPoint, rowIndex int, mapper MapperV2, columnName string) (interface{}, error) {
	mapperLen := len(mapper.Mappers)
	startPos := rowIndex * mapperLen
	if mapper.Multiple {
		return mapper.ParseStruct(columnName, excelPoints[startPos:startPos+mapperLen])
	}

	return mapper.ParseStruct("", excelPoints[startPos:startPos+mapperLen])
}
