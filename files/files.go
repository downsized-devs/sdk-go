// Package files contains small file-system helpers: extension parsing
// and existence checks that never panic on permission errors.
package files

import (
	"os"
	"strings"
)

// Extract file extension e.g. "txt, csv, docx" from full filename.
// If file has no extension then return empty string
func GetExtension(filename string) string {
	// This will not cause panic, because it will always return array with empty string,
	// even if the input string is empty
	filenamewithext := strings.Split(filename, ".")

	fileextension := filenamewithext[len(filenamewithext)-1]
	if fileextension == filename {
		return ""
	}
	return fileextension
}

// Checks if a file exists and is not a directory
func IsExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		// Treat any stat error (not-found, permission denied, etc.) as
		// "not a known regular file" so the caller never has to handle a
		// nil FileInfo dereference here.
		return false
	}

	return !info.IsDir()
}
