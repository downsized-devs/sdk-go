package parser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/downsized-devs/sdk-go/logger"
	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
	"go.uber.org/mock/gomock"
)

func Test_excelParser_UnmarshalTransform(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	ep := excelParser{
		log: logger,
	}
	f, _ := os.Open("./files/file-test.xlsx")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	g, _ := os.Open("./files/file-test-2.xlsx")
	buf2 := bytes.NewBuffer(nil)
	io.Copy(buf2, g)

	h, _ := os.Open("./files/file-test-empty.xlsx")
	buf3 := bytes.NewBuffer(nil)
	io.Copy(buf3, h)

	type args struct {
		ctx  context.Context
		blob []byte
		t    Transformer
		e    *excelize.File
		m    [][]ExcelPoint
	}
	tests := []struct {
		name    string
		args    args
		want    ExcelResult
		wantErr bool
		m       ExcelPoint
	}{
		{
			name: "unmarshal",
			args: args{
				ctx:  context.Background(),
				blob: []byte("whatever"),
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "unmarshal empty excel",
			args: args{
				ctx:  context.Background(),
				blob: buf3.Bytes(),
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "open reader success ",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
			},
			want: ExcelResult{
				Data:    []ExcelResultData{{}, {}, {}, {}},
				Details: []ExcelResultDetails{{}, {}, {}, {}},
			},
			wantErr: false,
		},
		{
			name: "slice error",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				e:    &excelize.File{},
				t: Transformer{
					SheetRegex: `^\/(?!\/)(.*?)`,
				},
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "filter error",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				e:    &excelize.File{},
				m: [][]ExcelPoint{
					{
						{
							Value:  "2.0",
							RowNum: 1,
							ColNum: 0,
						},
						{
							Value:  "bb",
							RowNum: 1,
							ColNum: 1,
						},
					},
					{
						{
							Value:  "3.0",
							RowNum: 2,
							ColNum: 0,
						},
						{
							Value:  "cc",
							RowNum: 2,
							ColNum: 1,
						},
					},
				},
				t: Transformer{
					SlicerNumCols: 3,
					Filters: []Filter{
						{
							Regex:  `^\/(?!\/)(.*?)`,
							ColNum: 1,
						},
						{
							Regex:  `^\/(?!\/)(.*?)`,
							ColNum: 1,
						},
					},
				},
			},
			want:    ExcelResult{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ep.UnmarshalTransform(tt.args.ctx, tt.args.blob, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("excelParser.UnmarshalTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTransformer_slice(t *testing.T) {
	tr := Transformer{}
	testFile, _ := os.Open("./files/file-test.xlsx")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, testFile)

	mockFile, _ := excelize.OpenReader(bytes.NewReader(buf.Bytes()))

	type args struct {
		ctx       context.Context
		f         *excelize.File
		sheetName []string
		t         Transformer
		s         string
	}
	tests := []struct {
		name            string
		args            args
		want            [][]ExcelPoint
		wantErr         bool
		SheetSeparators string
		SheetRegex      string
		SheetType       string
	}{
		{
			name: "slice",
			args: args{
				ctx: context.Background(),
				f:   mockFile,
			},
			want: [][]ExcelPoint{
				{},
				{
					{
						Value:  "2.00",
						RowNum: 1,
						ColNum: 0,
					},
					{
						Value:  "bb",
						RowNum: 1,
						ColNum: 1,
					},
				},
				{},
				{
					{
						Value:  "3.00",
						RowNum: 2,
						ColNum: 0,
					},
					{
						Value:  "cc",
						RowNum: 2,
						ColNum: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "slice error",
			args: args{
				ctx: context.Background(),
				f:   mockFile,
			},
			want:       nil,
			wantErr:    true,
			SheetRegex: `^\/(?!\/)(.*?)`,
		},
		{
			name: "sheet separator success",
			args: args{
				ctx:       context.Background(),
				f:         mockFile,
				sheetName: []string{},
				s:         "aaa",
				t: Transformer{
					SlicerInitialSkipCols: 1,
					SlicerNumCols:         3,
					SheetName:             "biji",
					SheetSeparators:       "aaa",
				},
			},
			want: [][]ExcelPoint{
				{
					{
						Value:    "Sheet1",
						RowNum:   0,
						ColNum:   0,
						Position: "",
					},
				},
				{
					{
						Value:    "Sheet1",
						RowNum:   0,
						ColNum:   0,
						Position: "",
					},
					{
						Value:    "bb",
						RowNum:   1,
						ColNum:   1,
						Position: "",
					},
					{
						Value:    "Sheet1",
						RowNum:   0,
						ColNum:   0,
						Position: "",
					},
				},
				{
					{
						Value:    "Sheet1",
						RowNum:   0,
						ColNum:   0,
						Position: "",
					},
				},
				{
					{
						Value:    "Sheet1",
						RowNum:   0,
						ColNum:   0,
						Position: "",
					},
					{
						Value:    "cc",
						RowNum:   2,
						ColNum:   1,
						Position: "",
					},
					{
						Value:    "Sheet1",
						RowNum:   0,
						ColNum:   0,
						Position: "",
					},
				},
			},
			wantErr:         false,
			SheetSeparators: "aaa",
		},
		{
			name: "appender",
			args: args{
				ctx:       context.Background(),
				f:         mockFile,
				sheetName: []string{},
				s:         "aaa",
				t: Transformer{
					SlicerInitialSkipCols: 1,
					SlicerNumCols:         3,
					SheetName:             "biji",
					SheetSeparators:       "",
				},
			},
			want:            [][]ExcelPoint{{{Value: "Sheet1", RowNum: 0, ColNum: 0, Position: ""}}, {{Value: "Sheet1", RowNum: 0, ColNum: 0, Position: ""}, {Value: "bb", RowNum: 1, ColNum: 1, Position: ""}, {Value: "Sheet1", RowNum: 0, ColNum: 0, Position: ""}}, {{Value: "Sheet1", RowNum: 0, ColNum: 0, Position: ""}}, {{Value: "Sheet1", RowNum: 0, ColNum: 0, Position: ""}, {Value: "cc", RowNum: 2, ColNum: 1, Position: ""}, {Value: "Sheet1", RowNum: 0, ColNum: 0, Position: ""}}},
			wantErr:         false,
			SheetSeparators: "",
			SheetType:       "dynamic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.SheetSeparators = tt.SheetSeparators
			tr.SheetRegex = tt.SheetRegex
			tr.SheetType = tt.SheetType
			got, err := tr.slice(tt.args.ctx, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transformer.slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTransformer_filter(t *testing.T) {
	tr := Transformer{}

	type args struct {
		ctx context.Context
		m   [][]ExcelPoint
	}
	tests := []struct {
		name    string
		args    args
		want    [][]ExcelPoint
		wantErr bool
		filter  []Filter
	}{
		{
			name: "filter",
			args: args{
				ctx: context.Background(),
				m:   [][]ExcelPoint{},
			},
			want:    [][]ExcelPoint{},
			wantErr: false,
		},
		{
			name: "filter regex success",
			args: args{
				ctx: context.Background(),
				m: [][]ExcelPoint{
					{
						{
							Value:  "123",
							ColNum: 0,
						},
					},
				},
			},
			want: [][]ExcelPoint{
				{
					{
						Value:  "123",
						ColNum: 0,
					},
				},
			},
			wantErr: false,
			filter: []Filter{
				{
					ColNum: 0,
					Regex:  "^[0-9]*$",
				},
			},
		},
		{
			name: "filter regex err",
			args: args{
				ctx: context.Background(),
				m: [][]ExcelPoint{
					{
						{
							Value:  "123",
							ColNum: 0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
			filter: []Filter{
				{
					ColNum: 0,
					Regex:  `^\/(?!\/)(.*?)`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.Filters = tt.filter
			got, err := tr.filter(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transformer.filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_filterEmptyRows(t *testing.T) {
	type args struct {
		ctx context.Context
		m   [][]ExcelPoint
	}
	tests := []struct {
		name    string
		args    args
		want    [][]ExcelPoint
		wantErr bool
	}{
		{
			name: "filter 2",
			args: args{
				ctx: context.Background(),
				m: [][]ExcelPoint{
					{
						{
							Value:  "123",
							ColNum: 0,
						},
					},
				},
			},
			want: [][]ExcelPoint{
				{
					{
						Value:  "123",
						ColNum: 0,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterEmptyRows(tt.args.ctx, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterEmptyRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmptyRow(t *testing.T) {
	type args struct {
		nCol int
	}
	tests := []struct {
		name    string
		args    args
		want    []ExcelPoint
		wantErr bool
	}{
		{
			name: "empty row",
			args: args{
				nCol: 1,
			},
			want: []ExcelPoint{
				{
					Value:  "",
					ColNum: 0,
					RowNum: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EmptyRow(tt.args.nCol); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmptyRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransformer_parse(t *testing.T) {
	tr := Transformer{}
	type args struct {
		ctx context.Context
		m   [][]ExcelPoint
	}
	tests := []struct {
		name    string
		args    args
		want    ExcelResult
		want1   []string
		wantErr bool
		mapper  []Mapper
	}{
		{
			name: "parse 2",
			args: args{
				ctx: context.Background(),
				m: [][]ExcelPoint{
					{
						{
							Value:    "123",
							ColNum:   0,
							RowNum:   0,
							Position: "A1",
						},
					},
				},
			},
			want: ExcelResult{
				Data: []ExcelResultData{
					{
						"2": "123",
					},
				},
				Details: []ExcelResultDetails{
					{
						"2": ExcelPoint{
							Value:    "123",
							ColNum:   0,
							RowNum:   0,
							Position: "A1",
						},
					},
				},
			},
			want1:   nil,
			wantErr: false,
			mapper: []Mapper{
				{
					Name: "2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr.Mappers = tt.mapper

			got, got1, err := tr.parse(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transformer.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transformer.parse() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Transformer.parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_cutMatrix(t *testing.T) {
	type args struct {
		ctx  context.Context
		m    [][]ExcelPoint
		nCol int
	}
	tests := []struct {
		name    string
		args    args
		want    [][]ExcelPoint
		wantErr bool
	}{
		{
			name: "cut matrix",
			args: args{
				ctx:  context.Background(),
				m:    [][]ExcelPoint{},
				nCol: 0,
			},
			want:    [][]ExcelPoint{},
			wantErr: false,
		},
		{
			name: "cut matrix 2",
			args: args{
				ctx: context.Background(),
				m: [][]ExcelPoint{
					{
						{
							Value:  "123",
							ColNum: 0,
						},
					},
				},
				nCol: 1,
			},
			want: [][]ExcelPoint{
				{
					{
						Value:  "123",
						ColNum: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cut matrix len > ncol",
			args: args{
				ctx: context.Background(),
				m: [][]ExcelPoint{
					{
						{
							Value:  "123",
							ColNum: 0,
							RowNum: 0,
						},
						{
							Value:  "1234",
							ColNum: 1,
							RowNum: 1,
						},
					},
				},
				nCol: 0,
			},
			want: [][]ExcelPoint{
				{},
				{
					{
						Value:  "123",
						ColNum: 0,
						RowNum: 0,
					},
					{
						Value:  "1234",
						ColNum: 1,
						RowNum: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "cut matrix len < ncol",
			args: args{
				ctx: context.Background(),
				m: [][]ExcelPoint{
					{
						{
							Value:  "123",
							ColNum: 0,
							RowNum: 0,
						},
						{
							Value:  "1234",
							ColNum: 1,
							RowNum: 1,
						},
					},
				},
				nCol: 4,
			},
			want: [][]ExcelPoint{
				{
					{
						Value:  "123",
						ColNum: 0,
						RowNum: 0,
					},
					{
						Value:  "1234",
						ColNum: 1,
						RowNum: 1,
					},
					{},
					{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutMatrix(tt.args.ctx, tt.args.m, tt.args.nCol); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cutMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapper_ParseString(t *testing.T) {
	m := Mapper{}

	var (
		mocktimeLayouts = []string{
			"2-Jan-2006",
			"2-Jan-06",
			"02-Jan-2006",
			"02-Jan-06",
			"01-02-06",
			"2006-01-02",
			"2006-01-2",
		}
	)

	var mockTime time.Time
	for _, mockLayoutVal := range mocktimeLayouts {
		mockTime, _ = time.Parse(mockLayoutVal, mockLayoutVal)
	}

	type args struct {
		input      string
		mockTimeLo []string
	}
	tests := []struct {
		name        string
		args        args
		want        interface{}
		wantErr     bool
		Type        string
		Relation    string
		RelationMap map[string]int64
	}{
		{
			name: "parse string err",
			args: args{
				input: "boo",
			},
			want:     0,
			wantErr:  true,
			Type:     "string",
			Relation: "whatever",
			RelationMap: map[string]int64{
				"whatever": 1,
			},
		},
		{
			name: "parse string succ",
			args: args{
				input: "boo",
			},
			want:     int64(1),
			wantErr:  false,
			Type:     "string",
			Relation: "whatever",
			RelationMap: map[string]int64{
				"boo": 1,
			},
		},
		{
			name: "parse string relation empty",
			args: args{
				input: "boo",
			},
			want:     "boo",
			wantErr:  false,
			Type:     "string",
			Relation: "",
			RelationMap: map[string]int64{
				"boo": 1,
			},
		},
		{
			name: "parse int",
			args: args{
				input: "1",
			},
			want:    int64(1),
			wantErr: false,
			Type:    "int",
		},
		{
			name: "parse float",
			args: args{
				input: "1.2",
			},
			want:    float64(1.2),
			wantErr: false,
			Type:    "float",
		},
		{
			name: "parse date",
			args: args{
				input: "2006-01-2",
			},
			want:    "2006-01-2",
			wantErr: false,
			Type:    "date",
		},
		{
			name: "parse float with space",
			args: args{
				input: " 1.2",
			},
			want:    float64(1.2),
			wantErr: false,
			Type:    "float",
		},
		{
			name: "parse dynamic date succ",
			args: args{
				input:      "2006-01-2",
				mockTimeLo: mocktimeLayouts,
			},
			want:    mockTime,
			wantErr: false,
			Type:    "dynamic_date",
		},
		{
			name: "parse dynamic date with space succ",
			args: args{
				input:      " 2006-01-2",
				mockTimeLo: mocktimeLayouts,
			},
			want:    mockTime,
			wantErr: false,
			Type:    "dynamic_date",
		},
		{
			name: "parse dynamic date err",
			args: args{
				input:      "2006ss-01-2",
				mockTimeLo: mocktimeLayouts,
			},
			want:    time.Time{},
			wantErr: false,
			Type:    "dynamic_date",
		},
		{
			name: "parse date:02-Jan-2006 succ",
			args: args{
				input:      "02-Jan-2006",
				mockTimeLo: mocktimeLayouts,
			},
			want:    mockTime,
			wantErr: false,
			Type:    "date:02-Jan-2006",
		},
		{
			name: "parse date:02-Jan-2006 err",
			args: args{
				input:      "02sss-Jan-2006",
				mockTimeLo: mocktimeLayouts,
			},
			want:    time.Time{},
			wantErr: true,
			Type:    "date:02-Jan-2006",
		},
		{
			name: "parse date:02-Jan-06 succ",
			args: args{
				input:      "02-Jan-06",
				mockTimeLo: mocktimeLayouts,
			},
			want:    mockTime,
			wantErr: false,
			Type:    "date:02-Jan-06",
		},
		{
			name: "parse date:02-Jan-06 succ 2",
			args: args{
				input:      "2-Jan-2006",
				mockTimeLo: mocktimeLayouts,
			},
			want:    mockTime,
			wantErr: false,
			Type:    "date:02-Jan-06",
		},
		{
			name: "parse date:02-Jan-06 err",
			args: args{
				input:      "2ss-Jan-2006",
				mockTimeLo: mocktimeLayouts,
			},
			want:    time.Time{},
			wantErr: true,
			Type:    "date:02-Jan-06",
		},
		{
			name: "parse date:01-02-06 succ",
			args: args{
				input:      "01-02-06",
				mockTimeLo: mocktimeLayouts,
			},
			want:    mockTime,
			wantErr: false,
			Type:    "date:01-02-06",
		},
		{
			name: "parse date:01-02-06 err",
			args: args{
				input:      "01s-02-06",
				mockTimeLo: mocktimeLayouts,
			},
			want:    time.Time{},
			wantErr: true,
			Type:    "date:01-02-06",
		},
		{
			name: "parse date:2006-01-02 succ",
			args: args{
				input:      "2006-01-02",
				mockTimeLo: mocktimeLayouts,
			},
			want:    mockTime,
			wantErr: false,
			Type:    "date:2006-01-02",
		},
		{
			name: "parse date:2006-01-02 err",
			args: args{
				input:      "20ss06-01-02",
				mockTimeLo: mocktimeLayouts,
			},
			want:    time.Time{},
			wantErr: true,
			Type:    "date:2006-01-02",
		},
	}
	for _, tt := range tests {
		m.Type = tt.Type
		m.Relation = tt.Relation
		m.RelationMap = tt.RelationMap
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.ParseString(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mapper.ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mapper.ParseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initExcel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	loggerMock := mock_log.NewMockInterface(ctrl)
	loggerMock.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	type args struct {
		log logger.Interface
	}
	tests := []struct {
		name string
		args args
		want ExcelInterface
	}{
		{
			name: "init success",
			args: args{
				log: loggerMock,
			},
			want: &excelParser{
				loggerMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initExcel(tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixMergedRow(t *testing.T) {
	testFile, _ := os.Open("./files/file-merged-test.xlsx")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, testFile)

	mockFile, _ := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	type args struct {
		ctx       context.Context
		f         *excelize.File
		sheetName string
		MergeCell excelize.MergeCell
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []string
	}{
		{
			name: "merge cell",
			args: args{
				ctx:       context.Background(),
				f:         mockFile,
				sheetName: "Sheet1",
			},
			wantErr: false,
		},
		{
			name: "merge succ",
			args: args{
				ctx:       context.Background(),
				f:         mockFile,
				sheetName: "Sheetss1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fixMergedRow(tt.args.ctx, tt.args.f, tt.args.sheetName); (err != nil) != tt.wantErr {
				t.Errorf("fixMergedRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_excelParser_UnmarshalWithAxis(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	ep := excelParser{
		log: logger,
	}
	f, _ := os.Open("./files/file-test.xlsx")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	h, _ := os.Open("./files/file-test-empty.xlsx")
	buf3 := bytes.NewBuffer(nil)
	io.Copy(buf3, h)

	type args struct {
		ctx  context.Context
		blob []byte
		t    Transformer
	}
	tests := []struct {
		name    string
		args    args
		want    ExcelResult
		wantErr bool
		m       ExcelPoint
	}{
		{
			name: "unmarshal",
			args: args{
				ctx:  context.Background(),
				blob: []byte("whatever"),
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "unmarshal empty excel",
			args: args{
				ctx:  context.Background(),
				blob: buf3.Bytes(),
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "invalid sheet name",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				t: Transformer{
					SheetRegex: "Sheet100",
				},
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "invalid cols count",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				t: Transformer{
					SheetRegex:    "Sheet1",
					SlicerNumCols: 4,
					SlicerNumRows: 3,
				},
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "invalid rows count",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				t: Transformer{
					SheetRegex:    "Sheet1",
					SlicerNumCols: 2,
					SlicerNumRows: 5,
				},
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "invalid axis",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				t: Transformer{
					SheetRegex:    "Sheet1",
					SlicerNumCols: 2,
					SlicerNumRows: 3,
					MappersV2: map[string]MapperV2{
						"test": {
							Name: "test",
							Mappers: []Mapper{
								{
									Axis: "",
								},
							},
						},
					},
				},
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "parse error",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				t: Transformer{
					SheetRegex:    "Sheet1",
					SlicerNumCols: 2,
					SlicerNumRows: 3,
					MappersV2: map[string]MapperV2{
						"test": {
							Name: "test",
							Mappers: []Mapper{
								{
									Type: "float",
									Axis: "B1",
								},
							},
						},
					},
				},
			},
			want:    ExcelResult{},
			wantErr: true,
		},
		{
			name: "all goods",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				t: Transformer{
					SheetRegex:    "Sheet1",
					SlicerNumCols: 2,
					SlicerNumRows: 3,
					MappersV2: map[string]MapperV2{
						"test": {
							Name: "test",
							Mappers: []Mapper{
								{
									Name: "Value",
									Type: "float",
									Axis: "A1",
								},
								{
									Name: "Title",
									Type: "string",
									Axis: "B1",
								},
							},
						},
					},
				},
			},
			want: ExcelResult{
				Data: []ExcelResultData{{"Value": 1.00, "Title": "aa"}},
				Details: []ExcelResultDetails{
					{
						"Value": ExcelPoint{Value: "1.00", ColNum: 1, RowNum: 1, Position: "A1"},
						"Title": ExcelPoint{Value: "aa", ColNum: 2, RowNum: 1, Position: "B1"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ep.UnmarshalWithAxis(tt.args.ctx, tt.args.blob, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("excelParser.UnmarshalWithAxis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_excelParser_UnmarshalTransformMultipleSheet(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	ep := excelParser{
		log: logger,
	}

	testFile, _ := os.Open("./files/file-test-production-plan-2.xlsx")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, testFile)

	mockTime, _ := time.Parse("02/01/2006", "31/01/2023")

	type args struct {
		ctx  context.Context
		blob []byte
		opt  ExcelOption
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]ExcelResult
		wantErr bool
	}{
		{
			name: "parser multiple sheet success",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				opt: ExcelOption{
					Type:    "productionplan",
					Version: "v2",
					Parsers: []Transformer{
						{
							Name:                  "stockingEstimation",
							SheetRegex:            "Estimasi Kebutuhan Benur",
							SheetName:             "Estimasi Kebutuhan Benur",
							SheetSeparators:       "",
							SheetType:             "single",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         7,
							MappersV2: map[string]MapperV2{
								"kolam": {
									Name:     "pondId",
									Type:     "string",
									Relation: "pond",
								},
								"luas": {
									Name:        "area",
									FilterRegex: "^[1-9][0-9]*$",
									Type:        "float",
								},
								"tanggaltebar": {
									Name: "stockingDate",
									Type: "dynamic_date",
								},
								"hatchery": {
									Name: "hatchery",
									Type: "string",
								},
								"postlarva(pl)": {
									Name: "postLarva",
									Type: "float",
								},
								"density/m²": {
									Name: "density",
									Type: "float",
								},
								"jumlahtebar": {
									Name: "stocking",
									Type: "int",
								},
							},
						},
						{
							Name:                  "blindFeedingProgram",
							SheetRegex:            "Blind Feeding Program",
							SheetName:             "Blind Feeding Program",
							SheetSeparators:       "",
							SheetType:             "single",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         17,
							MappersV2: map[string]MapperV2{
								"doc": {
									Name: "doc",
									Type: "int",
								},
								"tinggiair(m)": {
									Name: "waterHeight",
									Type: "float",
								},
								"luaskolam(m²)": {
									Name: "pondArea",
									Type: "float",
								},
								"density/m²": {
									Name: "density",
									Type: "float",
								},
								"stocking": {
									Name:        "stocking",
									Type:        "int",
									FilterRegex: "^[1-9][0-9]*$",
								},
								"pembagi": {
									Name: "divider",
									Type: "float",
								},
								"pkn/100rbekor": {
									Name: "pkn100",
									Type: "float",
								},
								"ajust/hr": {
									Name: "ajust",
									Type: "float",
								},
								"mbw": {
									Name: "mbw",
									Type: "float",
								},
								"adg": {
									Name: "adg",
									Type: "float",
								},
								"biomass": {
									Name: "biomass",
									Type: "float",
								},
								"sr": {
									Name: "sr",
									Type: "float",
								},
								"fcr": {
									Name: "fcr",
									Type: "float",
								},
								"p/h": {
									Name: "ph",
									Type: "float",
								},
								"feed.kum": {
									Name: "feedCum",
									Type: "float",
								},
								"nopakan": {
									Name: "feedNumber",
									Type: "float",
								},
							},
						},
						{
							Name:                  "prodPlanEstimatedFeed",
							SheetRegex:            "Estimasi Kebutuhan Pakan",
							SheetName:             "Estimasi Kebutuhan Pakan",
							SheetSeparators:       "",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         17,
							IsDynamicHeader:       true,
							MappersV2: map[string]MapperV2{
								"kolam": {
									Name:     "pondId",
									Type:     "string",
									Relation: "pond",
								},
								"pakankumulatif": {
									Name:        "feedCum",
									Type:        "float",
									FilterRegex: "^[1-9][0-9]*$",
								},
							},
						},
						{
							Name:                  "prodPlanFeedingProgramPerPond",
							SheetType:             "dynamic",
							SheetRegex:            "Feeding Program per Kola",
							SheetName:             "Indeks Feeding Program per Kola",
							SheetSeparators:       "",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         55,
							Filters: []Filter{
								{
									ColNum: 2,
									Regex:  "[0-9]+",
								},
							},
							Mappers: []Mapper{
								{
									Name:   "doc",
									ColNum: 0,
									Type:   "int",
								},
								{
									Name:   "stocking",
									ColNum: 1,
									Type:   "float",
								},
								{
									Name:   "divider",
									ColNum: 2,
									Type:   "float",
								},
								{
									Name:   "feedNumber",
									ColNum: 3,
									Type:   "float",
								},
								{
									Name:   "targetAbw",
									ColNum: 4,
									Type:   "float",
								},
								{
									Name:   "targetAdg",
									ColNum: 5,
									Type:   "float",
								},
								{
									Name:   "size",
									ColNum: 6,
									Type:   "float",
								},
								{
									Name:   "sr",
									ColNum: 7,
									Type:   "float",
								},
								{
									Name:   "estMort",
									ColNum: 8,
									Type:   "float",
								},
								{
									Name:   "population",
									ColNum: 9,
									Type:   "float",
								},
								{
									Name:   "biomass",
									ColNum: 10,
									Type:   "float",
								},
								{
									Name:   "index1",
									ColNum: 11,
									Type:   "float",
								},
								{
									Name:   "adg1",
									ColNum: 12,
									Type:   "float",
								},
								{
									Name:   "ph1",
									ColNum: 13,
									Type:   "float",
								},
								{
									Name:   "index2",
									ColNum: 14,
									Type:   "float",
								},
								{
									Name:   "adg2",
									ColNum: 15,
									Type:   "float",
								},
								{
									Name:   "ph2",
									ColNum: 16,
									Type:   "float",
								},
								{
									Name:   "index3",
									ColNum: 17,
									Type:   "float",
								},
								{
									Name:   "adg3",
									ColNum: 18,
									Type:   "float",
								},
								{
									Name:   "ph3",
									ColNum: 19,
									Type:   "float",
								},
								{
									Name:   "index4",
									ColNum: 20,
									Type:   "float",
								},
								{
									Name:   "adg4",
									ColNum: 21,
									Type:   "float",
								},
								{
									Name:   "ph4",
									ColNum: 22,
									Type:   "float",
								},
								{
									Name:   "index5",
									ColNum: 23,
									Type:   "float",
								},
								{
									Name:   "adg5",
									ColNum: 24,
									Type:   "float",
								},
								{
									Name:   "ph5",
									ColNum: 25,
									Type:   "float",
								},
								{
									Name:   "index6",
									ColNum: 26,
									Type:   "float",
								},
								{
									Name:   "adg6",
									ColNum: 27,
									Type:   "float",
								},
								{
									Name:   "ph6",
									ColNum: 28,
									Type:   "float",
								},
								{
									Name:   "index7",
									ColNum: 39,
									Type:   "float",
								},
								{
									Name:   "adg7",
									ColNum: 30,
									Type:   "float",
								},
								{
									Name:   "ph7",
									ColNum: 31,
									Type:   "float",
								},
								{
									Name:   "index8",
									ColNum: 32,
									Type:   "float",
								},
								{
									Name:   "adg8",
									ColNum: 33,
									Type:   "float",
								},
								{
									Name:   "ph8",
									ColNum: 34,
									Type:   "float",
								},
								{
									Name:   "index9",
									ColNum: 35,
									Type:   "float",
								},
								{
									Name:   "adg9",
									ColNum: 36,
									Type:   "float",
								},
								{
									Name:   "ph9",
									ColNum: 37,
									Type:   "float",
								},
								{
									Name:   "index10",
									ColNum: 38,
									Type:   "float",
								},
								{
									Name:   "adg10",
									ColNum: 39,
									Type:   "float",
								},
								{
									Name:   "ph10",
									ColNum: 40,
									Type:   "float",
								},
								{
									Name:   "area",
									ColNum: 42,
									Type:   "float",
								},
								{
									Name:   "waterHeight",
									ColNum: 43,
									Type:   "float",
								},
								{
									Name:   "chlorine",
									ColNum: 44,
									Type:   "float",
								},
								{
									Name:   "quickpro",
									ColNum: 45,
									Type:   "float",
								},
								{
									Name:   "thionnat",
									ColNum: 46,
									Type:   "float",
								},
								{
									Name:   "caOH2",
									ColNum: 47,
									Type:   "float",
								},
								{
									Name:   "caCO3",
									ColNum: 48,
									Type:   "float",
								},
								{
									Name:   "molasse2",
									ColNum: 49,
									Type:   "float",
								},
								{
									Name:   "feedCount",
									ColNum: 51,
									Type:   "float",
								},
								{
									Name:   "imunece5",
									ColNum: 52,
									Type:   "float",
								},
								{
									Name:   "vitaralAquatick",
									ColNum: 53,
									Type:   "float",
								},
								{
									Name:   "fastGrow5",
									ColNum: 54,
									Type:   "float",
								},
							},
						},
					},
				},
			},
			want: map[string]ExcelResult{
				"blindFeedingProgram": {
					Data: []ExcelResultData{
						{
							"doc":         int64(1),
							"waterHeight": float64(100),
							"pondArea":    float64(100),
							"density":     float64(100),
							"stocking":    int64(100),
							"divider":     float64(100),
							"pkn100":      float64(100),
							"ajust":       float64(100),
							"mbw":         float64(100),
							"adg":         float64(100),
							"biomass":     float64(100),
							"sr":          float64(100),
							"fcr":         float64(100),
							"ph":          float64(100),
							"feedCum":     float64(100),
							"feedNumber":  float64(100),
						},
					},
					Details: []ExcelResultDetails{
						{

							"doc":         ExcelPoint{Value: "1", RowNum: 3, ColNum: 0, Position: ""},
							"waterHeight": ExcelPoint{Value: "100", RowNum: 3, ColNum: 1, Position: ""},
							"pondArea":    ExcelPoint{Value: "100", RowNum: 3, ColNum: 2, Position: ""},
							"density":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 3, Position: ""},
							"stocking":    ExcelPoint{Value: "100", RowNum: 3, ColNum: 4, Position: ""},
							"divider":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 5, Position: ""},
							"pkn100":      ExcelPoint{Value: "100", RowNum: 3, ColNum: 6, Position: ""},
							"ajust":       ExcelPoint{Value: "100", RowNum: 3, ColNum: 7, Position: ""},
							"mbw":         ExcelPoint{Value: "100", RowNum: 3, ColNum: 8, Position: ""},
							"adg":         ExcelPoint{Value: "100", RowNum: 3, ColNum: 9, Position: ""},
							"biomass":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 10, Position: ""},
							"sr":          ExcelPoint{Value: "100", RowNum: 3, ColNum: 11, Position: ""},
							"fcr":         ExcelPoint{Value: "100", RowNum: 3, ColNum: 12, Position: ""},
							"ph":          ExcelPoint{Value: "100", RowNum: 3, ColNum: 13, Position: ""},
							"feedCum":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 14, Position: ""},
							"feedNumber":  ExcelPoint{Value: "100", RowNum: 3, ColNum: 15, Position: ""},
						},
					},
				},
				"stockingEstimation": {
					Data: []ExcelResultData{
						{
							"area":         float64(100),
							"density":      float64(100),
							"hatchery":     string("100"),
							"stocking":     int64(100),
							"pondId":       int(0),
							"postLarva":    float64(100),
							"stockingDate": mockTime,
						},
					},
					Details: []ExcelResultDetails{
						{
							"area":         ExcelPoint{Value: "100", RowNum: 4, ColNum: 1, Position: ""},
							"density":      ExcelPoint{Value: "100", RowNum: 4, ColNum: 5, Position: ""},
							"hatchery":     ExcelPoint{Value: "100", RowNum: 4, ColNum: 3, Position: ""},
							"stocking":     ExcelPoint{Value: "100", RowNum: 4, ColNum: 6, Position: ""},
							"pondId":       ExcelPoint{Value: "A1", RowNum: 4, ColNum: 0, Position: ""},
							"postLarva":    ExcelPoint{Value: "100", RowNum: 4, ColNum: 4, Position: ""},
							"stockingDate": ExcelPoint{Value: "31/01/2023", RowNum: 4, ColNum: 2, Position: ""},
						},
					},
				},
				"prodPlanEstimatedFeed": {
					Header: ExcelResultHeader{
						"pondId":                    ExcelHeader{Parent: "kolam", ParentValue: "Kolam", Name: "kolam", Value: "Kolam", ColIndex: 0},
						"merkpakan1:kodepakan1(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan1(kg)", Value: "Kode Pakan 1 (kg)", ColIndex: 1},
						"merkpakan1:kodepakan2(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan2(kg)", Value: "Kode Pakan 2 (kg)", ColIndex: 2},
						"merkpakan1:kodepakan3(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan3(kg)", Value: "Kode Pakan 3 (kg)", ColIndex: 3},
						"merkpakan1:kodepakan4(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan4(kg)", Value: "Kode Pakan 4 (kg)", ColIndex: 4},
						"merkpakan2:kodepakan1(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan1(kg)", Value: "Kode Pakan 1 (kg)", ColIndex: 5},
						"merkpakan2:kodepakan2(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan2(kg)", Value: "Kode Pakan 2 (kg)", ColIndex: 6},
						"merkpakan2:kodepakan3(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan3(kg)", Value: "Kode Pakan 3 (kg)", ColIndex: 7},
						"merkpakan2:kodepakan4(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan4(kg)", Value: "Kode Pakan 4 (kg)", ColIndex: 8},
						"merkpakan3:kodepakan1(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan1(kg)", Value: "Kode Pakan 1 (kg)", ColIndex: 9},
						"merkpakan3:kodepakan2(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan2(kg)", Value: "Kode Pakan 2 (kg)", ColIndex: 10},
						"merkpakan3:kodepakan3(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan3(kg)", Value: "Kode Pakan 3 (kg)", ColIndex: 11},
						"merkpakan3:kodepakan4(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan4(kg)", Value: "Kode Pakan 4 (kg)", ColIndex: 12},
						"feedCum":                   ExcelHeader{Parent: "pakankumulatif", ParentValue: "Pakan Kumulatif", Name: "pakankumulatif", Value: "Pakan Kumulatif", ColIndex: 13},
					},
					Data: []ExcelResultData{
						{
							"feedCum":                   float64(186),
							"merkpakan1:kodepakan1(kg)": string("10"),
							"merkpakan1:kodepakan2(kg)": string("11"),
							"merkpakan1:kodepakan3(kg)": string("12"),
							"merkpakan1:kodepakan4(kg)": string("13"),
							"merkpakan2:kodepakan1(kg)": string("14"),
							"merkpakan2:kodepakan2(kg)": string("15"),
							"merkpakan2:kodepakan3(kg)": string("16"),
							"merkpakan2:kodepakan4(kg)": string("17"),
							"merkpakan3:kodepakan1(kg)": string("18"),
							"merkpakan3:kodepakan2(kg)": string("19"),
							"merkpakan3:kodepakan3(kg)": string("20"),
							"merkpakan3:kodepakan4(kg)": string("21"),
							"pondId":                    int(0),
						},
					},
					Details: []ExcelResultDetails{
						{
							"feedCum":                   ExcelPoint{Value: "186", RowNum: 4, ColNum: 13, Position: ""},
							"merkpakan1:kodepakan1(kg)": ExcelPoint{Value: "10", RowNum: 4, ColNum: 1, Position: ""},
							"merkpakan1:kodepakan2(kg)": ExcelPoint{Value: "11", RowNum: 4, ColNum: 2, Position: ""},
							"merkpakan1:kodepakan3(kg)": ExcelPoint{Value: "12", RowNum: 4, ColNum: 3, Position: ""},
							"merkpakan1:kodepakan4(kg)": ExcelPoint{Value: "13", RowNum: 4, ColNum: 4, Position: ""},
							"merkpakan2:kodepakan1(kg)": ExcelPoint{Value: "14", RowNum: 4, ColNum: 5, Position: ""},
							"merkpakan2:kodepakan2(kg)": ExcelPoint{Value: "15", RowNum: 4, ColNum: 6, Position: ""},
							"merkpakan2:kodepakan3(kg)": ExcelPoint{Value: "16", RowNum: 4, ColNum: 7, Position: ""},
							"merkpakan2:kodepakan4(kg)": ExcelPoint{Value: "17", RowNum: 4, ColNum: 8, Position: ""},
							"merkpakan3:kodepakan1(kg)": ExcelPoint{Value: "18", RowNum: 4, ColNum: 9, Position: ""},
							"merkpakan3:kodepakan2(kg)": ExcelPoint{Value: "19", RowNum: 4, ColNum: 10, Position: ""},
							"merkpakan3:kodepakan3(kg)": ExcelPoint{Value: "20", RowNum: 4, ColNum: 11, Position: ""},
							"merkpakan3:kodepakan4(kg)": ExcelPoint{Value: "21", RowNum: 4, ColNum: 12, Position: ""},
							"pondId":                    ExcelPoint{Value: "A1", RowNum: 4, ColNum: 0, Position: ""},
						},
					},
				},
				"A1-prodPlanFeedingProgramPerPond": {
					Data: []ExcelResultData{
						{
							"divider":    float64(10),
							"doc":        int64(1),
							"feedNumber": float64(1),
							"stocking":   float64(100),
						},
					},
					Details: []ExcelResultDetails{
						{
							"divider":    ExcelPoint{Value: "10", RowNum: 3, ColNum: 2, Position: "C4"},
							"doc":        ExcelPoint{Value: "1", RowNum: 3, ColNum: 0, Position: "A4"},
							"feedNumber": ExcelPoint{Value: "1", RowNum: 3, ColNum: 3, Position: "D4"},
							"stocking":   ExcelPoint{Value: "100", RowNum: 3, ColNum: 1, Position: "B4"},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ep.UnmarshalTransformMultipleSheet(tt.args.ctx, tt.args.blob, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("excelParser.UnmarshalTransformMultipleSheet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("excelParser.UnmarshalTransformMultipleSheet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_excelParser_UnmarshalTransformMultipleSheetV2(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	ep := excelParser{
		log: logger,
	}

	testFile, _ := os.Open("./files/file-test-production-plan-2.xlsx")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, testFile)

	mockTime, _ := time.Parse("02/01/2006", "31/01/2023")

	type args struct {
		ctx  context.Context
		blob []byte
		opt  ExcelOption
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]ExcelResult
		wantErr bool
	}{
		{
			name: "parser multiple sheet success",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				opt: ExcelOption{
					Type:    "productionplan",
					Version: "v2",
					Parsers: []Transformer{
						{
							Name:                  "stockingEstimation",
							SheetRegex:            "Estimasi Kebutuhan Benur",
							SheetName:             "Estimasi Kebutuhan Benur",
							SheetSeparators:       "",
							SheetType:             "single",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         7,
							MappersV2: map[string]MapperV2{
								"kolam": {
									Name:     "pondId",
									Type:     "string",
									Relation: "pond",
								},
								"luas": {
									Name:        "area",
									FilterRegex: "^[1-9][0-9]*$",
									Type:        "float",
								},
								"tanggaltebar": {
									Name: "stockingDate",
									Type: "dynamic_date",
								},
								"hatchery": {
									Name: "hatchery",
									Type: "string",
								},
								"postlarva(pl)": {
									Name: "postLarva",
									Type: "float",
								},
								"density/m²": {
									Name: "density",
									Type: "float",
								},
								"jumlahtebar": {
									Name: "stocking",
									Type: "int",
								},
							},
						},
						{
							Name:                  "blindFeedingProgram",
							SheetRegex:            "Blind Feeding Program",
							SheetName:             "Blind Feeding Program",
							SheetSeparators:       "",
							SheetType:             "single",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         17,
							MappersV2: map[string]MapperV2{
								"doc": {
									Name: "doc",
									Type: "int",
								},
								"tinggiair(m)": {
									Name: "waterHeight",
									Type: "float",
								},
								"luaskolam(m²)": {
									Name: "pondArea",
									Type: "float",
								},
								"density/m²": {
									Name: "density",
									Type: "float",
								},
								"stocking": {
									Name:        "stocking",
									Type:        "int",
									FilterRegex: "^[1-9][0-9]*$",
								},
								"pembagi": {
									Name: "divider",
									Type: "float",
								},
								"pkn/100rbekor": {
									Name: "pkn100",
									Type: "float",
								},
								"ajust/hr": {
									Name: "ajust",
									Type: "float",
								},
								"mbw": {
									Name: "mbw",
									Type: "float",
								},
								"adg": {
									Name: "adg",
									Type: "float",
								},
								"biomass": {
									Name: "biomass",
									Type: "float",
								},
								"sr": {
									Name: "sr",
									Type: "float",
								},
								"fcr": {
									Name: "fcr",
									Type: "float",
								},
								"p/h": {
									Name: "ph",
									Type: "float",
								},
								"feed.kum": {
									Name: "feedCum",
									Type: "float",
								},
								"nopakan": {
									Name: "feedNumber",
									Type: "float",
								},
							},
						},
						{
							Name:                  "prodPlanEstimatedFeed",
							SheetRegex:            "Estimasi Kebutuhan Pakan",
							SheetName:             "Estimasi Kebutuhan Pakan",
							SheetSeparators:       "",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         17,
							IsDynamicHeader:       true,
							MappersV2: map[string]MapperV2{
								"kolam": {
									Name:     "pondId",
									Type:     "string",
									Relation: "pond",
								},
								"pakankumulatif": {
									Name:        "feedCum",
									Type:        "float",
									FilterRegex: "^[1-9][0-9]*$",
								},
							},
						},
						{
							Name:                  "prodPlanFeedingProgramPerPond",
							SheetType:             "dynamic",
							SheetRegex:            "Feeding Program per Kola",
							SheetName:             "Indeks Feeding Program per Kola",
							SheetSeparators:       "",
							SlicerInitialSkipRows: 2,
							SlicerInitialSkipCols: 0,
							SlicerNumRows:         0,
							SlicerNumCols:         55,
							Filters: []Filter{
								{
									ColNum: 2,
									Regex:  "[0-9]+",
								},
							},
							Mappers: []Mapper{
								{
									Name:   "doc",
									ColNum: 0,
									Type:   "int",
								},
								{
									Name:   "stocking",
									ColNum: 1,
									Type:   "float",
								},
								{
									Name:   "divider",
									ColNum: 2,
									Type:   "float",
								},
								{
									Name:   "feedNumber",
									ColNum: 3,
									Type:   "float",
								},
								{
									Name:   "targetAbw",
									ColNum: 4,
									Type:   "float",
								},
								{
									Name:   "targetAdg",
									ColNum: 5,
									Type:   "float",
								},
								{
									Name:   "size",
									ColNum: 6,
									Type:   "float",
								},
								{
									Name:   "sr",
									ColNum: 7,
									Type:   "float",
								},
								{
									Name:   "estMort",
									ColNum: 8,
									Type:   "float",
								},
								{
									Name:   "population",
									ColNum: 9,
									Type:   "float",
								},
								{
									Name:   "biomass",
									ColNum: 10,
									Type:   "float",
								},
								{
									Name:   "index1",
									ColNum: 11,
									Type:   "float",
								},
								{
									Name:   "adg1",
									ColNum: 12,
									Type:   "float",
								},
								{
									Name:   "ph1",
									ColNum: 13,
									Type:   "float",
								},
								{
									Name:   "index2",
									ColNum: 14,
									Type:   "float",
								},
								{
									Name:   "adg2",
									ColNum: 15,
									Type:   "float",
								},
								{
									Name:   "ph2",
									ColNum: 16,
									Type:   "float",
								},
								{
									Name:   "index3",
									ColNum: 17,
									Type:   "float",
								},
								{
									Name:   "adg3",
									ColNum: 18,
									Type:   "float",
								},
								{
									Name:   "ph3",
									ColNum: 19,
									Type:   "float",
								},
								{
									Name:   "index4",
									ColNum: 20,
									Type:   "float",
								},
								{
									Name:   "adg4",
									ColNum: 21,
									Type:   "float",
								},
								{
									Name:   "ph4",
									ColNum: 22,
									Type:   "float",
								},
								{
									Name:   "index5",
									ColNum: 23,
									Type:   "float",
								},
								{
									Name:   "adg5",
									ColNum: 24,
									Type:   "float",
								},
								{
									Name:   "ph5",
									ColNum: 25,
									Type:   "float",
								},
								{
									Name:   "index6",
									ColNum: 26,
									Type:   "float",
								},
								{
									Name:   "adg6",
									ColNum: 27,
									Type:   "float",
								},
								{
									Name:   "ph6",
									ColNum: 28,
									Type:   "float",
								},
								{
									Name:   "index7",
									ColNum: 39,
									Type:   "float",
								},
								{
									Name:   "adg7",
									ColNum: 30,
									Type:   "float",
								},
								{
									Name:   "ph7",
									ColNum: 31,
									Type:   "float",
								},
								{
									Name:   "index8",
									ColNum: 32,
									Type:   "float",
								},
								{
									Name:   "adg8",
									ColNum: 33,
									Type:   "float",
								},
								{
									Name:   "ph8",
									ColNum: 34,
									Type:   "float",
								},
								{
									Name:   "index9",
									ColNum: 35,
									Type:   "float",
								},
								{
									Name:   "adg9",
									ColNum: 36,
									Type:   "float",
								},
								{
									Name:   "ph9",
									ColNum: 37,
									Type:   "float",
								},
								{
									Name:   "index10",
									ColNum: 38,
									Type:   "float",
								},
								{
									Name:   "adg10",
									ColNum: 39,
									Type:   "float",
								},
								{
									Name:   "ph10",
									ColNum: 40,
									Type:   "float",
								},
								{
									Name:   "area",
									ColNum: 42,
									Type:   "float",
								},
								{
									Name:   "waterHeight",
									ColNum: 43,
									Type:   "float",
								},
								{
									Name:   "chlorine",
									ColNum: 44,
									Type:   "float",
								},
								{
									Name:   "quickpro",
									ColNum: 45,
									Type:   "float",
								},
								{
									Name:   "thionnat",
									ColNum: 46,
									Type:   "float",
								},
								{
									Name:   "caOH2",
									ColNum: 47,
									Type:   "float",
								},
								{
									Name:   "caCO3",
									ColNum: 48,
									Type:   "float",
								},
								{
									Name:   "molasse2",
									ColNum: 49,
									Type:   "float",
								},
								{
									Name:   "feedCount",
									ColNum: 51,
									Type:   "float",
								},
								{
									Name:   "imunece5",
									ColNum: 52,
									Type:   "float",
								},
								{
									Name:   "vitaralAquatick",
									ColNum: 53,
									Type:   "float",
								},
								{
									Name:   "fastGrow5",
									ColNum: 54,
									Type:   "float",
								},
							},
						},
					},
				},
			},
			want: map[string]ExcelResult{
				"blindFeedingProgram": {
					Data: []ExcelResultData{
						{
							"doc":         int64(1),
							"waterHeight": float64(100),
							"pondArea":    float64(100),
							"density":     float64(100),
							"stocking":    int64(100),
							"divider":     float64(100),
							"pkn100":      float64(100),
							"ajust":       float64(100),
							"mbw":         float64(100),
							"adg":         float64(100),
							"biomass":     float64(100),
							"sr":          float64(100),
							"fcr":         float64(100),
							"ph":          float64(100),
							"feedCum":     float64(100),
							"feedNumber":  float64(100),
						},
					},
					Details: []ExcelResultDetails{
						{

							"doc":         ExcelPoint{Value: "1", RowNum: 3, ColNum: 0, Position: ""},
							"waterHeight": ExcelPoint{Value: "100", RowNum: 3, ColNum: 1, Position: ""},
							"pondArea":    ExcelPoint{Value: "100", RowNum: 3, ColNum: 2, Position: ""},
							"density":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 3, Position: ""},
							"stocking":    ExcelPoint{Value: "100", RowNum: 3, ColNum: 4, Position: ""},
							"divider":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 5, Position: ""},
							"pkn100":      ExcelPoint{Value: "100", RowNum: 3, ColNum: 6, Position: ""},
							"ajust":       ExcelPoint{Value: "100", RowNum: 3, ColNum: 7, Position: ""},
							"mbw":         ExcelPoint{Value: "100", RowNum: 3, ColNum: 8, Position: ""},
							"adg":         ExcelPoint{Value: "100", RowNum: 3, ColNum: 9, Position: ""},
							"biomass":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 10, Position: ""},
							"sr":          ExcelPoint{Value: "100", RowNum: 3, ColNum: 11, Position: ""},
							"fcr":         ExcelPoint{Value: "100", RowNum: 3, ColNum: 12, Position: ""},
							"ph":          ExcelPoint{Value: "100", RowNum: 3, ColNum: 13, Position: ""},
							"feedCum":     ExcelPoint{Value: "100", RowNum: 3, ColNum: 14, Position: ""},
							"feedNumber":  ExcelPoint{Value: "100", RowNum: 3, ColNum: 15, Position: ""},
						},
					},
				},
				"stockingEstimation": {
					Data: []ExcelResultData{
						{
							"area":         float64(100),
							"density":      float64(100),
							"hatchery":     string("100"),
							"stocking":     int64(100),
							"pondId":       int(0),
							"postLarva":    float64(100),
							"stockingDate": mockTime,
						},
					},
					Details: []ExcelResultDetails{
						{
							"area":         ExcelPoint{Value: "100", RowNum: 4, ColNum: 1, Position: ""},
							"density":      ExcelPoint{Value: "100", RowNum: 4, ColNum: 5, Position: ""},
							"hatchery":     ExcelPoint{Value: "100", RowNum: 4, ColNum: 3, Position: ""},
							"stocking":     ExcelPoint{Value: "100", RowNum: 4, ColNum: 6, Position: ""},
							"pondId":       ExcelPoint{Value: "A1", RowNum: 4, ColNum: 0, Position: ""},
							"postLarva":    ExcelPoint{Value: "100", RowNum: 4, ColNum: 4, Position: ""},
							"stockingDate": ExcelPoint{Value: "31/01/2023", RowNum: 4, ColNum: 2, Position: ""},
						},
					},
				},
				"prodPlanEstimatedFeed": {
					Header: ExcelResultHeader{
						"pondId":                    ExcelHeader{Parent: "kolam", ParentValue: "Kolam", Name: "kolam", Value: "Kolam", ColIndex: 0},
						"merkpakan1:kodepakan1(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan1(kg)", Value: "Kode Pakan 1 (kg)", ColIndex: 1},
						"merkpakan1:kodepakan2(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan2(kg)", Value: "Kode Pakan 2 (kg)", ColIndex: 2},
						"merkpakan1:kodepakan3(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan3(kg)", Value: "Kode Pakan 3 (kg)", ColIndex: 3},
						"merkpakan1:kodepakan4(kg)": ExcelHeader{Parent: "merkpakan1", ParentValue: "Merk Pakan 1", Name: "kodepakan4(kg)", Value: "Kode Pakan 4 (kg)", ColIndex: 4},
						"merkpakan2:kodepakan1(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan1(kg)", Value: "Kode Pakan 1 (kg)", ColIndex: 5},
						"merkpakan2:kodepakan2(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan2(kg)", Value: "Kode Pakan 2 (kg)", ColIndex: 6},
						"merkpakan2:kodepakan3(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan3(kg)", Value: "Kode Pakan 3 (kg)", ColIndex: 7},
						"merkpakan2:kodepakan4(kg)": ExcelHeader{Parent: "merkpakan2", ParentValue: "Merk Pakan 2", Name: "kodepakan4(kg)", Value: "Kode Pakan 4 (kg)", ColIndex: 8},
						"merkpakan3:kodepakan1(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan1(kg)", Value: "Kode Pakan 1 (kg)", ColIndex: 9},
						"merkpakan3:kodepakan2(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan2(kg)", Value: "Kode Pakan 2 (kg)", ColIndex: 10},
						"merkpakan3:kodepakan3(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan3(kg)", Value: "Kode Pakan 3 (kg)", ColIndex: 11},
						"merkpakan3:kodepakan4(kg)": ExcelHeader{Parent: "merkpakan3", ParentValue: "Merk Pakan 3", Name: "kodepakan4(kg)", Value: "Kode Pakan 4 (kg)", ColIndex: 12},
						"feedCum":                   ExcelHeader{Parent: "pakankumulatif", ParentValue: "Pakan Kumulatif", Name: "pakankumulatif", Value: "Pakan Kumulatif", ColIndex: 13},
					},
					Data: []ExcelResultData{
						{
							"feedCum":                   float64(186),
							"merkpakan1:kodepakan1(kg)": string("10"),
							"merkpakan1:kodepakan2(kg)": string("11"),
							"merkpakan1:kodepakan3(kg)": string("12"),
							"merkpakan1:kodepakan4(kg)": string("13"),
							"merkpakan2:kodepakan1(kg)": string("14"),
							"merkpakan2:kodepakan2(kg)": string("15"),
							"merkpakan2:kodepakan3(kg)": string("16"),
							"merkpakan2:kodepakan4(kg)": string("17"),
							"merkpakan3:kodepakan1(kg)": string("18"),
							"merkpakan3:kodepakan2(kg)": string("19"),
							"merkpakan3:kodepakan3(kg)": string("20"),
							"merkpakan3:kodepakan4(kg)": string("21"),
							"pondId":                    int(0),
						},
					},
					Details: []ExcelResultDetails{
						{
							"feedCum":                   ExcelPoint{Value: "186", RowNum: 4, ColNum: 13, Position: ""},
							"merkpakan1:kodepakan1(kg)": ExcelPoint{Value: "10", RowNum: 4, ColNum: 1, Position: ""},
							"merkpakan1:kodepakan2(kg)": ExcelPoint{Value: "11", RowNum: 4, ColNum: 2, Position: ""},
							"merkpakan1:kodepakan3(kg)": ExcelPoint{Value: "12", RowNum: 4, ColNum: 3, Position: ""},
							"merkpakan1:kodepakan4(kg)": ExcelPoint{Value: "13", RowNum: 4, ColNum: 4, Position: ""},
							"merkpakan2:kodepakan1(kg)": ExcelPoint{Value: "14", RowNum: 4, ColNum: 5, Position: ""},
							"merkpakan2:kodepakan2(kg)": ExcelPoint{Value: "15", RowNum: 4, ColNum: 6, Position: ""},
							"merkpakan2:kodepakan3(kg)": ExcelPoint{Value: "16", RowNum: 4, ColNum: 7, Position: ""},
							"merkpakan2:kodepakan4(kg)": ExcelPoint{Value: "17", RowNum: 4, ColNum: 8, Position: ""},
							"merkpakan3:kodepakan1(kg)": ExcelPoint{Value: "18", RowNum: 4, ColNum: 9, Position: ""},
							"merkpakan3:kodepakan2(kg)": ExcelPoint{Value: "19", RowNum: 4, ColNum: 10, Position: ""},
							"merkpakan3:kodepakan3(kg)": ExcelPoint{Value: "20", RowNum: 4, ColNum: 11, Position: ""},
							"merkpakan3:kodepakan4(kg)": ExcelPoint{Value: "21", RowNum: 4, ColNum: 12, Position: ""},
							"pondId":                    ExcelPoint{Value: "A1", RowNum: 4, ColNum: 0, Position: ""},
						},
					},
				},
				"A1-prodPlanFeedingProgramPerPond": {
					Data: []ExcelResultData{
						{
							"divider":    float64(10),
							"doc":        int64(1),
							"feedNumber": float64(1),
							"stocking":   float64(100),
						},
					},
					Details: []ExcelResultDetails{
						{
							"divider":    ExcelPoint{Value: "10", RowNum: 3, ColNum: 2, Position: "C4"},
							"doc":        ExcelPoint{Value: "1", RowNum: 3, ColNum: 0, Position: "A4"},
							"feedNumber": ExcelPoint{Value: "1", RowNum: 3, ColNum: 3, Position: "D4"},
							"stocking":   ExcelPoint{Value: "100", RowNum: 3, ColNum: 1, Position: "B4"},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ep.UnmarshalTransformMultipleSheetV2(tt.args.ctx, tt.args.blob, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("excelParser.UnmarshalTransformMultipleSheetV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("excelParser.UnmarshalTransformMultipleSheetV2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_excelParser_Marshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	p := &excelParser{
		log: logger,
	}

	type args struct {
		ctx        context.Context
		input      []ExcelInput
		pathToFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: []ExcelInput{
					{
						SheetName: "Register",
						SheetValue: map[string]string{
							"A1": "Facility Name",
							"B1": "Sample Name",
							"C1": "Sample Number",
							"D1": "Address",
							"A2": "Devs",
							"B2": "A1",
							"C2": "123",
							"D2": "Jl. Damai Indah Sentosa No 420, RT 6 RW 9",
						},
					},
					{
						SheetName: "Result",
						SheetValue: map[string]string{
							"A1": "Facility Name",
							"B1": "Sample Name",
							"C1": "Status",
							"A2": "Devs",
							"B2": "A1",
							"C2": "Positive",
							"A3": "Devs",
							"B3": "A2",
							"C3": "Negative",
						},
					},
				},
				pathToFile: "./files/file-test-generated.xlsx",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.Marshal(tt.args.ctx, tt.args.input, tt.args.pathToFile); (err != nil) != tt.wantErr {
				t.Errorf("excelParser.Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			e := os.Remove("./files/file-test-generated.xlsx")
			if e != nil {
				t.Fatal(e)
			}
		})
	}
}

func Test_excelParser_WriteReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	parser := &excelParser{
		log: logger,
	}

	mockNotValidBlob := &bytes.Buffer{}

	mockExcelFile := excelize.NewFile()
	mockExcelFile.SetCellValue("Sheet1", "A1", "Hello")
	mockExcelFile.SetCellValue("Sheet1", "B1", "world!")
	mockBlob, _ := mockExcelFile.WriteToBuffer()

	mockWantOk := make(map[string]string)
	mockWantOk["A1"] = "Hello"
	mockWantOk["B1"] = "world!"

	type args struct {
		ctx      context.Context
		blob     []byte
		ioReader io.Writer
	}

	tests := []struct {
		name      string
		args      args
		sheetName string
		want      map[string]string
		wantErr   bool
	}{
		{
			name: "error invalid blob",
			args: args{
				ctx:      context.Background(),
				blob:     mockNotValidBlob.Bytes(),
				ioReader: &bytes.Buffer{},
			},
			sheetName: "Sheet1",
			want:      mockWantOk,
			wantErr:   true,
		},
		{
			name: "all ok",
			args: args{
				ctx:      context.Background(),
				blob:     mockBlob.Bytes(),
				ioReader: &bytes.Buffer{},
			},
			sheetName: "Sheet1",
			want:      mockWantOk,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.WriteReader(tt.args.ctx, tt.args.blob, tt.args.ioReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("excelParser.WriteReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check per value
			if err == nil {
				f2, err := excelize.OpenReader(bytes.NewReader(tt.args.blob))
				assert.NoError(t, err)

				for key, value := range tt.want {
					cellValue, err := f2.GetCellValue(tt.sheetName, key)
					assert.NoError(t, err)
					assert.Equal(t, fmt.Sprint(value), cellValue)
				}
			}
		})
	}

}

func Test_excelParser_InjectValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	parser := excelParser{
		log: logger,
	}

	mockNotValidBlob := &bytes.Buffer{}

	mockFileOk := excelize.NewFile()
	mockFileOk.SetCellValue("Sheet1", "A1", "A1")

	bufferFileOK, _ := mockFileOk.WriteToBuffer()

	mockInjectValueOk := make(map[string]interface{})
	mockInjectValueOk["B1"] = "B1"
	mockInjectValueOk["C1"] = 3.14
	mockInjectValueOk["D1"] = 42

	type args struct {
		ctx         context.Context
		blobByte    *bytes.Buffer
		sheetName   string
		injectValue map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error not valid blob byte",
			args: args{
				sheetName:   "Sheet1",
				blobByte:    mockNotValidBlob,
				injectValue: mockInjectValueOk,
			},
			wantErr: true,
		},
		{
			name: "all ok 1",
			args: args{
				sheetName:   "Sheet1",
				blobByte:    bufferFileOK,
				injectValue: mockInjectValueOk,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injectedBuffer, err := parser.InjectValue(tt.args.ctx, tt.args.blobByte.Bytes(), tt.args.sheetName, tt.args.injectValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("excelParser.InjectValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check per value
			if err == nil {
				f2, err := excelize.OpenReader(bytes.NewReader(injectedBuffer))
				assert.NoError(t, err)

				for key, value := range tt.args.injectValue {
					cellValue, err := f2.GetCellValue(tt.args.sheetName, key)
					assert.NoError(t, err)
					assert.Equal(t, fmt.Sprint(value), cellValue)
				}
			}
		})
	}
}

func Test_excelParser_UnmarshalTransformYAxis(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	ep := excelParser{
		log: logger,
	}
	f, _ := os.Open("./files/file-test-financial.xlsx")
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	type args struct {
		ctx  context.Context
		blob []byte
		t    Transformer
	}
	tests := []struct {
		name    string
		args    args
		want    ExcelResult
		wantErr bool
		m       ExcelPoint
	}{
		{
			name: "success",
			args: args{
				ctx:  context.Background(),
				blob: buf.Bytes(),
				t: Transformer{
					Name:                  "financialProfitLoss",
					SheetType:             "single",
					SheetRegex:            "Template",
					SheetName:             "Template",
					SheetSeparators:       "",
					SlicerInitialSkipRows: 3,
					SlicerInitialSkipCols: 0,
					SlicerNumRows:         0,
					SlicerNumCols:         0,
					IsDynamicHeader:       true,
					MappersV2: map[string]MapperV2{
						"farmname": {
							Name: "name",
							Type: "string",
						},
						"unit": {
							Name: "unit",
							Type: "string",
						},
						"historical": {
							MultiColumn: true,
							Multiple:    true,
							Name:        "historical",
							Mappers: []Mapper{
								{
									Name: "cycle",
									Type: "string",
								},
							},
						},
					},
					MapperYAxis: MapperYAxis{
						KeysColumn: []string{"farmname", "unit"},
						MappersVal: map[string]Mapper{
							"Beginning Date-Date": {
								Name: "startDate",
								Type: "dynamic_date",
							},
						},
					},
				},
			},
			wantErr: false,
			want: ExcelResult{
				Data: []ExcelResultData{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ep.UnmarshalTransformYAxis(tt.args.ctx, tt.args.blob, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("excelParser.UnmarshalTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Header, got.Header)
		})
	}
}
