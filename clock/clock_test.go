package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_timelib_GetCurrentTime(t *testing.T) {
	tl := Init()

	mockTime := time.Date(2021, 6, 1, 16, 38, 04, 0, time.UTC)
	Now = func() time.Time {
		return mockTime
	}
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "ok",
			want: mockTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tl.GetCurrentTime(); got != tt.want {
				t.Errorf("timelib.GetCurrentTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timelib_SubstractTime(t *testing.T) {
	tl := Init()

	mockOrigin := time.Date(2021, 6, 1, 16, 38, 04, 0, time.UTC)
	mockDeductor := time.Date(2021, 6, 1, 16, 38, 03, 0, time.UTC)

	type args struct {
		origin   time.Time
		deductor time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "ok",
			args: args{
				origin:   mockOrigin,
				deductor: mockDeductor,
			},
			want: time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tl.SubstractTime(tt.args.origin, tt.args.deductor); got != tt.want {
				t.Errorf("timelib.SubstractTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timelib_AddTime(t *testing.T) {
	tl := Init()

	mockOrigin := time.Date(2021, 6, 1, 16, 38, 04, 0, time.UTC)
	mockAdder := time.Second

	type args struct {
		origin time.Time
		adder  time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "ok",
			args: args{
				origin: mockOrigin,
				adder:  mockAdder,
			},
			want: time.Date(2021, 6, 1, 16, 38, 05, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tl.AddTime(tt.args.origin, tt.args.adder); got != tt.want {
				t.Errorf("timelib.AddTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timelib_GetTimeInLocation(t *testing.T) {
	tl := Init()

	mockLocationParam := AsiaJakarta
	location, err := time.LoadLocation(string(mockLocationParam))
	if err != nil {
		panic(err)
	}

	mockTimeParam := time.Date(2021, 6, 1, 16, 38, 04, 0, time.UTC)

	type args struct {
		locationParam Location
		timeParam     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				locationParam: mockLocationParam,
				timeParam:     mockTimeParam,
			},
			want:    time.Date(2021, 6, 1, 23, 38, 04, 0, location),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tl.GetTimeInLocation(tt.args.locationParam, tt.args.timeParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("timelib.GetTimeInLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("timelib.GetTimeInLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timelib_GetFirstDayOfTheMonth(t *testing.T) {
	tl := Init()

	type args struct {
		year  int
		month time.Month
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "ok",
			args: args{
				year:  2021,
				month: 6,
			},
			want: time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tl.GetFirstDayOfTheMonth(tt.args.year, tt.args.month); got != tt.want {
				t.Errorf("timelib.GetFirstDayOfTheMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timelib_GetLastDayOfTheMonth(t *testing.T) {
	tl := Init()

	type args struct {
		year  int
		month time.Month
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "ok",
			args: args{
				year:  2021,
				month: 6,
			},
			want: time.Date(2021, 6, 30, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tl.GetLastDayOfTheMonth(tt.args.year, tt.args.month); got != tt.want {
				t.Errorf("timelib.GetLastDayOfTheMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timelib_ConvertFromString(t *testing.T) {
	tl := Init()

	type args struct {
		timeFormat string
		timeString string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				timeFormat: "2006-01-02 15:04:05",
				timeString: "2021-06-01 16:38:04",
			},
			want:    time.Date(2021, 6, 1, 16, 38, 04, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tl.ConvertFromString(tt.args.timeFormat, tt.args.timeString)
			if (err != nil) != tt.wantErr {
				t.Errorf("timelib.ConvertFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_timelib_ConvertToString(t *testing.T) {
	tl := Init()

	type args struct {
		timeFormat string
		timeParam  time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				timeFormat: "2006-01-02 15:04:05",
				timeParam:  time.Date(2021, 6, 1, 16, 38, 04, 0, time.UTC),
			},
			want: "2021-06-01 16:38:04",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tl.ConvertToString(tt.args.timeFormat, tt.args.timeParam); got != tt.want {
				t.Errorf("timelib.ConvertToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
