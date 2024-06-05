package num

import "math"

// Rounding given number to n decimal place
func RoundingNumber(num float64, n int) float64 {
	base := math.Pow(10, float64(n))

	return math.Round(num*base) / base
}
