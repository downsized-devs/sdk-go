package timelib

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

type timelib struct {
}

func Init() Interface {
	return &timelib{}
}

func (t *timelib) GetCurrentTime() time.Time {
	return Now()
}

func (t *timelib) SubstractTime(origin, deductor time.Time) time.Duration {
	return origin.Sub(deductor)
}

func (t *timelib) AddTime(origin time.Time, adder time.Duration) time.Time {
	return origin.Add(adder)
}

func (t *timelib) GetTimeInLocation(locationParam Location, timeParam time.Time) (time.Time, error) {
	location, err := time.LoadLocation(string(locationParam))
	if err != nil {
		return time.Time{}, err
	}

	return timeParam.In(location), nil
}

func (t *timelib) GetFirstDayOfTheMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func (t *timelib) GetLastDayOfTheMonth(year int, month time.Month) time.Time {
	firstDayOfNextMonth := t.GetFirstDayOfTheMonth(year, month+1)
	return firstDayOfNextMonth.AddDate(0, 0, -1)
}

func (t *timelib) ConvertFromString(timeFormat, timeString string) (time.Time, error) {
	return time.Parse(timeFormat, timeString)
}

func (t *timelib) ConvertToString(timeFormat string, timeParam time.Time) string {
	return timeParam.Format(timeFormat)
}
