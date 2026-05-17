// Package num contains numeric utility helpers — currently a safe
// division primitive that hides divide-by-zero handling behind a
// configurable fallback.
package num

// SafeDivision returns numerator/denominator. When denominator is zero
// the function returns 0 if zeroValue is true; otherwise it returns the
// numerator unchanged so callers can choose an identity value.
func SafeDivision(numerator float64, denominator float64, zeroValue bool) float64 {
	if denominator == 0 {
		if !zeroValue {
			return numerator
		} else {
			return 0
		}
	} else {
		return numerator / denominator
	}
}
