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
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		info += fmt.Sprintf("In function %s()", funcName)
	}
	for i := 0; ok; i++ {
		info += fmt.Sprintf("\n%s:%d", file, line)
		_, file, line, ok = runtime.Caller(skip + i)
	}
	return info
}

func PrintCallStack() {
	debug.PrintStack()
}

func LogCallStack() {
	log.Print(Stack())
}

func Stack() string {
	return string(debug.Stack())
}

func formatValue(value interface{}) string {
	return fmt.Sprintf("\n     Type: %T\n    Value: %v\nGo Syntax: %#v", value, value, value)
}

func formatCallstack(skip int) string {
	return fmt.Sprintf("\nCallstack: %s", CallStackInfo(skip + 1))
}

func FormatSkip(skip int, value interface{}) string {
	return formatValue(value) + formatCallstack(skip + 1)
}

func Format(value interface{}) string {
	return FormatSkip(2, value)
}

func Dump(values ...interface{}) {
	if Logger != nil {
		for _, value := range values {
			Logger.Println(formatValue(value))
		}
	}
	Logger.Println(formatCallstack(2))
}

func Print(values ...interface{}) {
	if Logger != nil {
		Logger.Print(values...)
	}
}

func Printf(format string, values ...interface{}) {
	if Logger != nil {
		Logger.Printf(format, values...)
	}
}

func Quit() {
	os.Exit(-1)
}
