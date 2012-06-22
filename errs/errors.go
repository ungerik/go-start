// Errors and error utility functions.
package errs

import (
	"bytes"
	"fmt"
	"github.com/ungerik/go-start/debug"
)

///////////////////////////////////////////////////////////////////////////////
// ErrNotImplemented

type ErrNotImplemented string

func (self ErrNotImplemented) Error() string {
	if self == "" {
		return "Not implemented"
	}
	return string(self) + " not implemented"
}

///////////////////////////////////////////////////////////////////////////////
// ErrIndexOutOfBounds

type ErrIndexOutOfBounds struct {
	What   string
	Index  int
	Length int
}

func (self *ErrIndexOutOfBounds) Error() string {
	max := self.Length - 1
	return fmt.Sprintf("%s index %d out of bounds [0..%d]", self.What, self.Index, max)
}

func IfErrIndexOutOfBounds(what string, index int, length int) *ErrIndexOutOfBounds {
	if index < 0 || index >= length {
		if Config.FormatWithCallStack {
			what += debug.CallStackInfo(2)
		}
		return &ErrIndexOutOfBounds{what, index, length}
	}
	return nil
}

func PanicIfErrIndexOutOfBounds(what string, index int, length int) {
	if index < 0 || index >= length {
		if Config.FormatWithCallStack {
			what += debug.CallStackInfo(2)
		}
		panic(&ErrIndexOutOfBounds{what, index, length})
	}
}

///////////////////////////////////////////////////////////////////////////////
// ErrSlice

// ErrSlice implements error for a slice of errors.
// Error() returns the Error() results for every slice field
// concaternated by '\n'. nil errors are ignored.
type ErrSlice []error

func (self ErrSlice) Error() string {
	return self.Join("\n")
}

func (self ErrSlice) Join(sep string) string {
	var buf bytes.Buffer
	for _, err := range self {
		if err != nil {
			if buf.Len() > 0 {
				buf.WriteString(sep)
			}
			buf.WriteString(err.Error())
		}
	}
	return buf.String()
}

func Errors(errs ...error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}
	// for i := 0; i < len(errs); i++ {
	// 	if multi, ok := errs[i].(ErrSlice); ok {
	// 		errs = append(append(errs[0:i], multi...), errs[i+1:]...)
	// 		i += len(multi)
	// 	}
	// }
	return ErrSlice(errs)
}
