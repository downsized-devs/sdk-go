package num

func EmptyStringSlice(length int) []string {
	result := []string{}
	for i := 0; i < length; i++ {
		result = append(result, "")
	}
	return result
}
