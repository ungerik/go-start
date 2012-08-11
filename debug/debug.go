// Debug functions.
package debug

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
)

// Nop is a dummy function that can be called in source files where
// other debug functions are constantly added and removed.
// That way import "github.com/ungerik/go-start/debug" won't cause an error when
// no other debug function is currently used.
// Arbitrary objects can be passed as arguments to avoid "declared and not used"
// error messages when commenting code out and in.
// The result is a nil interface{} dummy value.
func Nop(dummiesIn ...interface{}) (dummyOut interface{}) {
	return nil
}

// Set Logger to nil to suppress debug output
var Logger = log.New(os.Stdout, "", 0)

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

func PrintCallStack() {
	fmt.Println(string(debug.Stack()))
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
