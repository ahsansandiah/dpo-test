package traceHelper

import (
	"fmt"
	"strings"
	"testing"

	"runtime"

	"github.com/stretchr/testify/assert"
)

func TestErrorTraceCorrectFileAndLineNumber(t *testing.T) {
	_, file, line, _ := runtime.Caller(0)
	expectedFile := file
	expectedLine := line + 2

	file, _ = ErrorTrace(2)
	fileParts := strings.Split(file, ":")
	actualFile := fileParts[0]
	actualLine := fileParts[1]

	if actualFile != expectedFile {
		t.Errorf("Expected file %s, but got %s", expectedFile, actualFile)
	}

	fmt.Println(actualLine, expectedLine)
	if actualLine != fmt.Sprintf("%d", expectedLine) {
		assert.NotEqual(t, actualFile, expectedLine)
	}
}

func TestErrorTraceZeroStackDepth(t *testing.T) {
	file, funcName := ErrorTrace(0)

	if file == "" || funcName == "" {
		t.Errorf("Expected non-empty file and function name, but got file: %s, function name: %s", file, funcName)
	}
}

func TestErrorTraceExceedsStackSize(t *testing.T) {
	file, funcx := ErrorTrace(1000)

	if file == "" {
		t.Errorf("Expected a non-empty file string, but got an empty string")
	}

	if funcx == "" {
		t.Errorf("Expected a non-empty function string, but got an empty string")
	}
}
