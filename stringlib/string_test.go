package stringlib

import (
	"math/rand"
	"strings"
	"testing"
	"time"
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

// TestRandStringBytesRepeatability checks if setting the same seed results in the same output
func TestRandStringBytesRepeatability(t *testing.T) {
	length := 10
	seed := time.Now().UnixNano()

	rand.New(rand.NewSource(seed))
	str1 := RandStringBytes(length)

	rand.New(rand.NewSource(seed))
	str2 := RandStringBytes(length)

	if str1 != str2 {
		t.Errorf("Expected the same string output for the same seed, but got '%s' and '%s'", str1, str2)
	}
}
