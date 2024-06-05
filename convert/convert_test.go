package convert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToArrInt64(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []int64
		wantErr bool
	}{
		{
			name:    "error to []int64",
			args:    args{i: fmt.Errorf("huehue")},
			want:    []int64{},
			wantErr: true,
		},
		{
			name:    "[]error to []int64",
			args:    args{i: []error{fmt.Errorf("huehue")}},
			want:    []int64{},
			wantErr: true,
		},
		{
			name:    "int64 to []int64",
			args:    args{i: int64(5)},
			want:    []int64{5},
			wantErr: false,
		},
		{
			name:    "[]string to []int64",
			args:    args{i: []string{"1", "2", "3"}},
			want:    []int64{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "[]float64 to []int64",
			args:    args{i: []float64{1.43, 2.3, 3.9}},
			want:    []int64{1, 2, 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToArrInt64(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToArrInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToArrInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToChar(t *testing.T) {
	tests := []struct {
		name string
		args int
		want string
	}{
		{
			name: "1 to A",
			args: 1,
			want: "A",
		},
		{
			name: "3 to C",
			args: 3,
			want: "C",
		},
		{
			name: "25 to Y",
			args: 25,
			want: "Y",
		},
		{
			name: "26 to Z",
			args: 26,
			want: "Z",
		},
		{
			name: "0 to empty string",
			args: 0,
			want: "",
		},
		{
			name: "-1 to empty string",
			args: -1,
			want: "",
		},
		{
			name: "27 to AA",
			args: 27,
			want: "AA",
		},
		{
			name: "51 to AY",
			args: 51,
			want: "AY",
		},
		{
			name: "52 to AZ",
			args: 52,
			want: "AZ",
		},
		{
			name: "53 to BA",
			args: 53,
			want: "BA",
		},
		{
			name: "702 to ZZ",
			args: 702,
			want: "ZZ",
		},
		{
			name: "703 to AAA",
			args: 703,
			want: "AAA",
		},
		{
			name: "703 to BXQ",
			args: 1993,
			want: "BXQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToChar(tt.args); got != tt.want {
				t.Errorf("IntToChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPascalCaseToCamelCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "Test empty string",
			args: "",
			want: "",
		},
		{
			name: "Pascal to Camel",
			args: "PasCal",
			want: "pasCal",
		},
		{
			name: "Camel to Camel",
			args: "caMel",
			want: "caMel",
		},
		{
			name: "Test single uppercase string",
			args: "A",
			want: "a",
		},
		{
			name: "Test single lowercase string",
			args: "a",
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PascalCaseToCamelCase(tt.args); got != tt.want {
				t.Errorf("PascalCaseToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	flo := float64(123.43)
	pointerFlo := &flo
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "error",
			args:    args{i: fmt.Errorf("huehue")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "[]error",
			args:    args{i: []error{fmt.Errorf("huehue")}},
			want:    1, // print the length
			wantErr: false,
		},
		{
			name:    "random string",
			args:    args{i: "huehue"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "[]float64",
			args:    args{i: []float64{5, 6, 7}},
			want:    3, // print the length
			wantErr: false,
		},
		{
			name:    "int64",
			args:    args{i: int64(12)},
			want:    12,
			wantErr: false,
		},
		{
			name:    "pointer to float",
			args:    args{i: flo},
			want:    123.43,
			wantErr: false,
		},
		{
			name:    "float dereference",
			args:    args{i: *pointerFlo},
			want:    123.43,
			wantErr: false,
		},
		{
			name:    "string float",
			args:    args{i: "123.43"},
			want:    123.43,
			wantErr: false,
		},
		{
			name:    "nil input",
			args:    args{i: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "error to int64",
			args:    args{i: fmt.Errorf("ayonima")},
			want:    int64(0),
			wantErr: true,
		},
		{
			name:    "[]int64 to int64",
			args:    args{i: []int64{1, 2, 3}},
			want:    int64(3),
			wantErr: false,
		},
		{
			name:    "string to int64",
			args:    args{i: "69420"},
			want:    69420,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestToString(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "convert int to string",
			args:    args{i: int64(2)},
			want:    "2",
			wantErr: false,
		},
		{
			name:    "convert array to string",
			args:    args{i: []int64{1, 2, 3}},
			want:    "[1 2 3]",
			wantErr: false,
		},
		{
			name:    "convert from float to string",
			args:    args{i: float64(6.9)},
			want:    "6.9",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToString(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRomanToInt64(t *testing.T) {
	type args struct {
		roman string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "case 1",
			args: args{roman: "MCMXCIV"},
			want: 1994,
		},
		{
			name: "case 2",
			args: args{roman: "LVIII"},
			want: 58,
		},
		{
			name: "case 3",
			args: args{roman: "IX"},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RomanToInt64(tt.args.roman); got != tt.want {
				t.Errorf("RomanToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64ToRoman(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{num: 1994},
			want: "MCMXCIV",
		},
		{
			name: "case 2",
			args: args{num: 58},
			want: "LVIII",
		},
		{
			name: "case 3",
			args: args{num: 9},
			want: "IX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64ToRoman(tt.args.num); got != tt.want {
				t.Errorf("IntToRoman() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToLowerCamelCase(t *testing.T) {
	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				input: "Additional Item",
			},
			want: "additionalItem",
		},
		{
			name: "case 2",
			args: args{
				input: "Application 1",
			},
			want: "application1",
		},
		{
			name: "case 3",
			args: args{
				input: "Green Algae (Chlorophyta)",
			},
			want: "greenAlgaeChlorophyta",
		},
		{
			name: "case 4",
			args: args{
				input: "TOM (Total Organic Matter)",
			},
			want: "tOMTotalOrganicMatter",
		},
		{
			name: "case 5",
			args: args{
				input: "Rasio N:P",
			},
			want: "rasioNP",
		},
		{
			name: "case 6",
			args: args{
				input: "Nitrite (NO₂)",
			},
			want: "nitriteNO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamelCase(tt.args.input); got != tt.want {
				t.Errorf("ToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				input: "Additional Item",
			},
			want: "AdditionalItem",
		},
		{
			name: "case 2",
			args: args{
				input: "Application 1",
			},
			want: "Application1",
		},
		{
			name: "case 3",
			args: args{
				input: "Green Algae (Chlorophyta)",
			},
			want: "GreenAlgaeChlorophyta",
		},
		{
			name: "case 4",
			args: args{
				input: "TOM (Total Organic Matter)",
			},
			want: "TOMTotalOrganicMatter",
		},
		{
			name: "case 5",
			args: args{
				input: "Rasio N:P",
			},
			want: "RasioNP",
		},
		{
			name: "case 6",
			args: args{
				input: "Nitrite (NO₂)",
			},
			want: "NitriteNO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPascalCase(tt.args.input); got != tt.want {
				t.Errorf("ToPascalCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
