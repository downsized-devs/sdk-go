package clock

import (
	"time"
)

type Location string

const (
	// Location
	AsiaJakarta   Location = "Asia/Jakarta"   // GMT+7
	AsiaSingapore Location = "Asia/Singapore" // GMT+8
	AsiaJayapura  Location = "Asia/Jayapura"  // GMT+9
)

var Now = time.Now

type Interface interface {
	GetCurrentTime() time.Time
	SubstractTime(origin, deductor time.Time) time.Duration
	AddTime(origin time.Time, adder time.Duration) time.Time
	GetTimeInLocation(locationParam Location, timeParam time.Time) (time.Time, error)
	GetFirstDayOfTheMonth(year int, month time.Month) time.Time
	GetLastDayOfTheMonth(year int, month time.Month) time.Time
	ConvertFromString(timeFormat, timeString string) (time.Time, error)
	ConvertToString(timeFormat string, timeParam time.Time) string
}

type clock struct {
}

func Init() Interface {
	return &clock{}
}

func (t *clock) GetCurrentTime() time.Time {
	return Now()
}

func (t *clock) SubstractTime(origin, deductor time.Time) time.Duration {
	return origin.Sub(deductor)
}

func (t *clock) AddTime(origin time.Time, adder time.Duration) time.Time {
	return origin.Add(adder)
}

func (t *clock) GetTimeInLocation(locationParam Location, timeParam time.Time) (time.Time, error) {
	location, err := time.LoadLocation(string(locationParam))
	if err != nil {
		return time.Time{}, err
	}

	return timeParam.In(location), nil
}

func (t *clock) GetFirstDayOfTheMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func (t *clock) GetLastDayOfTheMonth(year int, month time.Month) time.Time {
	firstDayOfNextMonth := t.GetFirstDayOfTheMonth(year, month+1)
	return firstDayOfNextMonth.AddDate(0, 0, -1)
}

func (t *clock) ConvertFromString(timeFormat, timeString string) (time.Time, error) {
	return time.Parse(timeFormat, timeString)
}

func (t *clock) ConvertToString(timeFormat string, timeParam time.Time) string {
	return timeParam.Format(timeFormat)
}
