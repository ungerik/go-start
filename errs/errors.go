// Errors and error utility functions.
package errs

import (
	"bytes"
	"fmt"
	"github.com/ungerik/go-start/debug"
)

///////////////////////////////////////////////////////////////////////////////
// NotImplemented

type NotImplemented string

func (self NotImplemented) Error() string {
	if self == "" {
		return "Not implemented"
	}
	return string(self) + " not implemented"
}

///////////////////////////////////////////////////////////////////////////////
// IndexOutOfBounds

type IndexOutOfBounds struct {
	What   string
	Index  int
	Length int
}

func (self *IndexOutOfBounds) Error() string {
	max := self.Length - 1
	return fmt.Sprintf("%s index %d out of bounds [0..%d]", self.What, self.Index, max)
}

func IfIndexOutOfBounds(what string, index int, length int) *IndexOutOfBounds {
	if index < 0 || index >= length {
		if Config.FormatWithCallStack {
			what += debug.CallStackInfo(2)
		}
		return &IndexOutOfBounds{what, index, length}
	}
	return nil
}

func PanicIfIndexOutOfBounds(what string, index int, length int) {
	if index < 0 || index >= length {
		if Config.FormatWithCallStack {
			what += debug.CallStackInfo(2)
		}
		panic(&IndexOutOfBounds{what, index, length})
	}
}

///////////////////////////////////////////////////////////////////////////////
// MultipleErrors

// MultipleErrors implements error for a slice of errors.
// Error() returns the Error() results for every slice field
// concaternated by '\n'. nil errors are ignored.
type MultipleErrors []error

func (self MultipleErrors) Error() string {
	return self.Join("\n")
}

func (self MultipleErrors) Join(sep string) string {
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
	// 	if multi, ok := errs[i].(MultipleErrors); ok {
	// 		errs = append(append(errs[0:i], multi...), errs[i+1:]...)
	// 		i += len(multi)
	// 	}
	// }
	return MultipleErrors(errs)
}
