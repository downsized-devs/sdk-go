package helper

import (
	"regexp"
	"strings"
)

// This function is converting text from camelCase to snake_case {textExample => test_example}
func ConvertToSnakeCase(text string) string {
	regex := regexp.MustCompile(`([A-Z])`)
	replaced := regex.ReplaceAllString(text, "_$1")
	result := strings.ToLower(replaced)
	return result
}

// This function is converting text from camelCase to UpperCase {textExample => Test Example}
func ConvertToUpperCase(text string) string {
	texts := strings.Split(text, "")
	texts[0] = strings.ToUpper(texts[0])
	result := strings.Join(texts, "")
	return result
}

// This function is to add space before every UpperCase character {textExample => test Example}
func convertToSpace(text string) string {
	regex := regexp.MustCompile(`([A-Z])`)
	result := regex.ReplaceAllString(text, " $1")
	return result
}

// This function is converting text from camelCase to Upper Case With Space {textExample => Text Example}
func ConvertToUpperSpace(text string) string {
	result := convertToSpace(text)
	result = ConvertToUpperCase(result)
	return result
}

// This function is converting text from camelCase to lower case with space {textExample => test example}
func ConvertToLowerSpace(text string) string {
	result := convertToSpace(text)
	result = strings.ToLower(result)
	return result
}

// This function is converting text from UpperCase to camelCase {TextExample => textExample}
func ConvertToCamelCase(text string) string {
	texts := strings.Split(text, "")
	texts[0] = strings.ToLower(texts[0])
	result := strings.Join(texts, "")
	return result
}

// This function is converting text from camelCase to lower-case-with-dash {textExample => text-example}
func ConvertToLowerDash(text string) string {
	regex := regexp.MustCompile(`([A-Z])`)
	replaced := regex.ReplaceAllString(text, "-$1")
	result := strings.ToLower(replaced)
	return result
}
