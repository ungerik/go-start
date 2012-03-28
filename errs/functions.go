package errs

import (
	"fmt"
	"github.com/ungerik/go-start/debug"
)

// FormatSkipStackFrames formats an os.Error with call stack information if FormatWithCallStack is true.
// It skips skip stack frames.
func FormatSkipStackFrames(skip int, format string, args ...interface{}) error {
	if Config.FormatWithCallStack {
		format += "\n" + debug.CallStackInfo(skip+1)
	}
	return fmt.Errorf(format, args...)
}

// Format formats an os.Error with call stack information if FormatWithCallStack is true
func Format(format string, args ...interface{}) error {
	return FormatSkipStackFrames(2, format, args...)
}

func Assert(condition bool, description string, args ...interface{}) {
	if !condition {
		panic(FormatSkipStackFrames(2, "Failed assertion: "+description, args...))
	}
}

// todo test
//func Catch(err *os.Error) {
//	if err == nil {
//		panic(FormatSkipStackFrames(2, "errs.Catch expects non nil os.Error pointer"))
//	}
//	CatchCall(func(e os.Error) {
//		*err = e
//	})
//}
//
//func CatchCall(callback func(os.Error)) {
//	if r := recover(); r != nil {
//		var err os.Error
//		if e, ok := r.(runtime.Error); ok {
//			panic(e)
//			//*err = e
//		} else if e, ok := r.(os.Error); ok {
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

// Panic if any of the args is a non nil os.Error
// Be careful because everything with a String() method qualifies as an os.Error
func AsPanic(args ...interface{}) {
	for _, arg := range args {
		if arg != nil {
			if err, ok := arg.(error); ok {
				panic(err)
			}
		}
	}
}

// Panic if the last element of args is a non nil os.Error
// Be careful because everything with a String() method qualifies as an os.Error
func LastAsPanic(args ...interface{}) {
	arg := args[len(args)-1]
	if arg != nil {
		if err, ok := arg.(error); ok {
			panic(err)
		}
	}
}

// Returns the first os.Error in args
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
