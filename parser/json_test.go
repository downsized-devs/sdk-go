package parser

import (
	"reflect"
	"testing"

	mock_log "github.com/downsized-devs/sdk-go/tests/mock/logger"
	"go.uber.org/mock/gomock"
)

type MockFarm struct {
	FarmId         int    `json:"farmId"`
	CompanyId      int    `json:"companyId"`
	FarmName       string `json:"farmName,omitempty"`
	FarmCode       string `json:"farmCode"`
	FarmLocation   string `json:"farmLocation"`
	FarmCoordinate string `json:"farmCoordinate"`
	PondCount      int    `json:"pondCount"`
}

func Test_jsonParser_MarshalWithSchemaValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()

	jsonSchema := initJSON(JSONOptions{
		Schema: map[string]string{
			"test_schema": "file://./files/test.schema.json",
		},
	}, logger)

	mockFarm := MockFarm{
		FarmId:         1,
		CompanyId:      1,
		FarmName:       "test",
		FarmCode:       "test-CODE",
		FarmLocation:   "test-location",
		FarmCoordinate: "{longitude},{latitude}",
		PondCount:      1,
	}
	mockFarmMarshalled, _ := jsonSchema.Marshal(mockFarm)

	mockFarmTestError := MockFarm{
		FarmId:         2,
		CompanyId:      2,
		FarmCode:       "test-CODE",
		FarmLocation:   "test-location",
		FarmCoordinate: "{longitude},{latitude}",
		PondCount:      2,
	}

	type args struct {
		sch  string
		orig interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success marshal with schema validation",
			args: args{
				sch:  "test_schema",
				orig: mockFarm,
			},
			want:    mockFarmMarshalled,
			wantErr: false,
		},
		{
			name: "failed marshal with schema validation: schema not found",
			args: args{
				sch:  "test_schema_not_found",
				orig: mockFarm,
			},
			wantErr: true,
		},
		{
			name: "failed marshal with schema validation: farmName is required",
			args: args{
				sch:  "test_schema",
				orig: mockFarmTestError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonSchema.MarshalWithSchemaValidation(tt.args.sch, tt.args.orig)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonParser.MarshalWithSchemaValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonParser.MarshalWithSchemaValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonParser_UnmarshalWithSchemaValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()

	jsonSchema := initJSON(JSONOptions{
		Schema: map[string]string{
			"test_schema": "file://./files/test.schema.json",
		},
	}, logger)

	type args struct {
		sch  string
		blob []byte
		dest interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success unmarshal with schema validation",
			args: args{
				sch: "test_schema",
				blob: []byte(`{
					"farmId":1,
					"companyId":1,
					"farmName":"test",
					"farmCode":"test-CODE",
					"farmLocation":"test-location",
					"farmCoordinate":"{longitude},{latitude}",
					"pondCount":1
				 }`),
				dest: &MockFarm{},
			},
			wantErr: false,
		},
		{
			name: "failed unmarshal with schema validation: schema not found",
			args: args{
				sch: "test_schema_not_found",
				blob: []byte(`{
					"farmId":1,
					"companyId":1,
					"farmName":"test",
					"farmCode":"test-CODE",
					"farmLocation":"test-location",
					"farmCoordinate":"{longitude},{latitude}",
					"pondCount":1
				 }`),
				dest: &MockFarm{},
			},
			wantErr: true,
		},
		{
			name: "failed unmarshal with schema validation: farmName is required",
			args: args{
				sch: "test_schema",
				blob: []byte(`{
					"farmId":1,
					"companyId":1,
					"farmCode":"test-CODE",
					"farmLocation":"test-location",
					"farmCoordinate":"{longitude},{latitude}",
					"pondCount":1
				 }`),
				dest: &MockFarm{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := jsonSchema.UnmarshalWithSchemaValidation(tt.args.sch, tt.args.blob, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("jsonParser.UnmarshalWithSchemaValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
