package utils

import (
	"fmt"
	"io"
	"strconv"
)

type StringWriter struct {
	writer io.Writer
}

func (self *StringWriter) Write(strings ...string) *StringWriter {
	for _, str := range strings {
		self.writer.Write([]byte(str))
	}
	return self
}

func (self *StringWriter) Printf(format string, args ...interface{}) *StringWriter {
	fmt.Fprintf(self.writer, format, args...)
	return self
}

func (self *StringWriter) Byte(value byte) *StringWriter {
	self.writer.Write([]byte{value})
	return self
}

func (self *StringWriter) WriteBytes(bytes []byte) *StringWriter {
	self.writer.Write(bytes)
	return self
}

func (self *StringWriter) Int(value int) *StringWriter {
	self.writer.Write([]byte(strconv.Itoa(value)))
	return self
}

func (self *StringWriter) Uint(value uint) *StringWriter {
	self.writer.Write([]byte(strconv.FormatUint(uint64(value), 10)))
	return self
}

func (self *StringWriter) Float(value float64) *StringWriter {
	self.writer.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
	return self
}

func (self *StringWriter) Bool(value bool) *StringWriter {
	self.writer.Write([]byte(strconv.FormatBool(value)))
	return self
}
