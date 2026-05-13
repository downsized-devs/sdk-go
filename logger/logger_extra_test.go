package logger

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_logger_Fatal_Exits verifies that Fatal terminates the process by
// re-invoking the test binary as a subprocess. The subprocess is detected
// via the LOGGER_FATAL_CRASHER env var.
func Test_logger_Fatal_Exits(t *testing.T) {
	if os.Getenv("LOGGER_FATAL_CRASHER") == "1" {
		l := Init(Config{Level: "debug"})
		l.Fatal(context.Background(), "intentional fatal exit for test")
		// Should be unreachable.
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^Test_logger_Fatal_Exits$")
	cmd.Env = append(os.Environ(), "LOGGER_FATAL_CRASHER=1")
	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	assert.True(t, ok, "Fatal should exit non-zero, got err=%v", err)
	if ok {
		assert.False(t, exitErr.Success())
	}
}

// Test_logger_Init_InvalidLevel verifies that Init aborts when given an
// unparseable log level. Same subprocess trick as above.
func Test_logger_Init_InvalidLevel(t *testing.T) {
	if os.Getenv("LOGGER_INIT_CRASHER") == "1" {
		_ = Init(Config{Level: "not-a-real-level"})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^Test_logger_Init_InvalidLevel$")
	cmd.Env = append(os.Environ(), "LOGGER_INIT_CRASHER=1")
	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	assert.True(t, ok, "Init with invalid level should exit non-zero, got err=%v", err)
	if ok {
		assert.False(t, exitErr.Success())
	}
}

func Test_logger_Init_CustomSkipFrameCount(t *testing.T) {
	// Just ensures the override branch is exercised; we don't assert on
	// the resulting frame because zerolog doesn't expose it cleanly.
	l := Init(Config{Level: "info", CallerSkipFrameCount: 7})
	assert.NotNil(t, l)
}

func Test_DefaultLogger(t *testing.T) {
	assert.NotNil(t, DefaultLogger())
}
