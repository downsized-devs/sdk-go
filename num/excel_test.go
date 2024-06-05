package num

import (
	"context"
	"reflect"
	"testing"
)

func TestExcelGenerateCoords(t *testing.T) {
	type args struct {
		ctx             context.Context
		coordinateRange string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "start column invalid",
			args: args{
				ctx:             context.Background(),
				coordinateRange: "210:D13",
			},
			wantErr: true,
			want:    []string{},
		},
		{
			name: "end column invalid",
			args: args{
				ctx:             context.Background(),
				coordinateRange: "D10:813",
			},
			wantErr: true,
			want:    []string{},
		},
		{
			name: "generate coords success",
			args: args{
				ctx:             context.Background(),
				coordinateRange: "D10:D13",
			},
			wantErr: false,
			want:    []string{"D10", "D11", "D12", "D13"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExcelGenerateCoords(tt.args.ctx, tt.args.coordinateRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExcelGenerateCoords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExcelGenerateCoords() = %v, want %v", got, tt.want)
			}
		})
	}
}
