package num

// / Checks denominator in a fraction (numerator/denominator) which returns zero value or original value depends on the zeroValue
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
