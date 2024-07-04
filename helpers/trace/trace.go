package traceHelper

import (
	"fmt"
	"runtime"
	"strings"
)

func ErrorTrace(deep int) (string, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(deep, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	file := fmt.Sprintf("%s:%d", frame.File, frame.Line)

	funcx := frame.Function
	funcxSplit := strings.Split(funcx, "/")
	funcx = funcxSplit[len(funcxSplit)-1]
	funcxSplit = strings.Split(funcx, ".")
	funcx = funcxSplit[len(funcxSplit)-1] + "()"

	return file, funcx
}
