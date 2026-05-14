package errors

import (
	"fmt"

	"github.com/downsized-devs/sdk-go/codes"
)

type stacktrace struct {
	message  string
	cause    error
	code     codes.Code
	file     string
	function string
	line     int
}

func (st *stacktrace) Error() string {
	return fmt.Sprint(st.message)
}

func (st *stacktrace) ExitCode() int {
	if st.code == codes.NoCode {
		return 1
	}
	return int(st.code)
}

// Unwrap returns the underlying cause so that errors.Is and errors.As can
// traverse the full error chain.
func (st *stacktrace) Unwrap() error {
	return st.cause
}
