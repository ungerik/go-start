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
	var buf bytes.Buffer
	for _, err := range self {
		if err != nil {
			if buf.Len() > 0 {
				buf.WriteByte('\n')
			}
			buf.WriteString(err.Error())
		}
	}
	return buf.String()
}
