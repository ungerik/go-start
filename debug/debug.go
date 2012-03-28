// Debug functions.
package debug

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// Set Logger to nil to suppress debug output
var Logger *log.Logger = log.New(os.Stdout, "", 0)

func CallStackInfo(skip int) (info string) {
	if pc, file, line, ok := runtime.Caller(skip); ok {
		funcName := runtime.FuncForPC(pc).Name()
		info += fmt.Sprintf("%s()\n%s:%d", funcName, file, line)
	}
	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		info += fmt.Sprintf("\n%s:%d", file, line)
	}
	return
}

func FormatSkip(skip int, value interface{}) string {
	callstack := CallStackInfo(skip + 1)
	return fmt.Sprintf("\n     Type: %T\n    Value: %v\nGo Syntax: %#v\nCallstack: %s", value, value, value, callstack)
}

func Format(value interface{}) string {
	return FormatSkip(2, value)
}

func Dump(values ...interface{}) {
	if Logger == nil {
		return
	}
	for _, value := range values {
		Logger.Println(FormatSkip(2, value))
	}
}

func Print(values ...interface{}) {
	if Logger == nil {
		return
	}
	Logger.Print(values...)
}

func Quit() {
	os.Exit(-1)
}

// Nop is a dummy function that can be called in source files where
// other debug functions are constantly added and removed.
// That way import "github.com/ungerik/go-start/debug" won't cause an error when
// no other debug function is called.
func Nop() {
}
