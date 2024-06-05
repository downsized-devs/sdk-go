package character

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CapitalizeFirstCharacter(s string, lan language.Tag) string {
	// Split the string into words
	words := strings.Fields(s)

	caser := cases.Title(lan)

	// Capitalize the first letter of each word
	for i, word := range words {
		// Capitalize the first letter of the word
		words[i] = caser.String(word)
	}

	// Join the words back into a string
	result := strings.Join(words, " ")

	return result
}

func IsStrongCharCombination(character string, minLength int) bool {
	// Define password strength criteria
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString

	specialCharRegex := `[!@#$%^&*(),.?":{}|<>+'=_\-\\;+~\[\]` + "`]"
	hasSpecialChar := regexp.MustCompile(specialCharRegex).MatchString

	// Check each criteria
	if len(character) < minLength {
		return false
	}
	if !hasUppercase(character) {
		return false
	}
	if !hasLowercase(character) {
		return false
	}
	if !hasDigit(character) {
		return false
	}
	if !hasSpecialChar(character) {
		return false
	}

	return true
}
