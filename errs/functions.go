package errs

import (
	"fmt"
	"github.com/ungerik/go-start/debug"
)

// FormatSkipStackFrames formats an error with call stack information if FormatWithCallStack is true.
// It skips skip stack frames.
func FormatSkipStackFrames(skip int, format string, args ...interface{}) error {
	if Config.FormatWithCallStack {
		format += "\n" + debug.CallStackInfo(skip+1)
	}
	return fmt.Errorf(format, args...)
}

// Format formats an error with call stack information if FormatWithCallStack is true
func Format(format string, args ...interface{}) error {
	return FormatSkipStackFrames(2, format, args...)
}

func Assert(condition bool, description string, args ...interface{}) {
	if !condition {
		panic(FormatSkipStackFrames(2, "Failed assertion: "+description, args...))
	}
}

// todo test
//func Catch(err *error) {
//	if err == nil {
//		panic(FormatSkipStackFrames(2, "errs.Catch expects non nil error pointer"))
//	}
//	CatchCall(func(e error) {
//		*err = e
//	})
//}
//
//func CatchCall(callback func(error)) {
//	if r := recover(); r != nil {
//		var err error
//		if e, ok := r.(runtime.Error); ok {
//			panic(e)
//			//*err = e
//		} else if e, ok := r.(error); ok {
//			err = e
//		} else if s, ok := r.(fmt.Stringer); ok {
//			err = os.NewError(s.String())
//		} else if s, ok := r.(string); ok {
//			err = os.NewError(s)
//		} else {
//			panic(r)
//		}
//		callback(err)
//	}
//}

// Panic if any of the args is a non nil error
func PanicOnError(args ...interface{}) {
	for _, arg := range args {
		if arg != nil {
			if err, ok := arg.(error); ok {
				panic(err)
			}
		}
	}
}

// Panic if the last element of args is a non nil error
func LastPanicOnError(args ...interface{}) {
	arg := args[len(args)-1]
	if arg != nil {
		if err, ok := arg.(error); ok {
			panic(err)
		}
	}
}

// Returns the first error in args
func First(args ...interface{}) error {
	for _, arg := range args {
		if arg != nil {
			if err, ok := arg.(error); ok {
				return err
			}
		}
	}
	return nil
}
