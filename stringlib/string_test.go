package stringlib

import (
	"strings"
	"testing"
)

// TestRandStringBytesLength tests if the output string has the expected length
func TestRandStringBytesLength(t *testing.T) {
	length := 10
	str := RandStringBytes(length)
	if len(str) != length {
		t.Errorf("Expected string of length %d, but got %d", length, len(str))
	}
}

// TestRandStringBytesContent checks if the output string contains only the allowed characters
func TestRandStringBytesContent(t *testing.T) {
	length := 10
	str := RandStringBytes(length)
	for _, char := range str {
		if !contains(letterBytes, byte(char)) {
			t.Errorf("The character '%c' is not allowed", char)
		}
	}
}

// contains checks if a byte is within a string
func contains(s string, b byte) bool {
	return string(b) == s || strings.ContainsRune(s, rune(b))
}
