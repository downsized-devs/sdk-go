package dates

import "time"

func Difference(a, b time.Time) int64 {
	if a.After(b) {
		a, b = b, a
	}
	return int64(b.Sub(a).Hours() / 24)
}
