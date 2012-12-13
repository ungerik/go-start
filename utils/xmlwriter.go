package utils

import (
	"fmt"
	"html"
	"io"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/reflection"
	//	"github.com/ungerik/go-start/debug"
)

///////////////////////////////////////////////////////////////////////////////
// XMLWriter

func NewXMLWriter(writer io.Writer) *XMLWriter {
	if xmlWriter, ok := writer.(*XMLWriter); ok {
		return xmlWriter
	}
	return &XMLWriter{writer: writer}
}

type XMLWriter struct {
	writer    io.Writer
	tagStack  []string
	inOpenTag bool
}

func (self *XMLWriter) WriteXMLDeclaration() *XMLWriter {
	return self.Content(`<?xml version="1.0" encoding="UTF-8"?>`)
}

func (self *XMLWriter) OpenTag(tag string) *XMLWriter {
	self.finishOpenTag()

	self.writer.Write([]byte{'<'})
	self.writer.Write([]byte(tag))

	self.tagStack = append(self.tagStack, tag)
	self.inOpenTag = true

	return self
}

// value will be HTML escaped and concaternated
func (self *XMLWriter) Attrib(name string, value ...interface{}) *XMLWriter {
	errs.Assert(self.inOpenTag, "utils.XMLWriter.Attrib() must be called inside of open tag")

	fmt.Fprintf(self.writer, " %s='", name)
	for _, valuePart := range value {
		str := html.EscapeString(fmt.Sprint(valuePart))
		self.writer.Write([]byte(str))
	}
	self.writer.Write([]byte{'\''})

	return self
}

func (self *XMLWriter) AttribIfNotDefault(name string, value interface{}) *XMLWriter {
	if reflection.IsDefaultValue(value) {
		return self
	}
	return self.Attrib(name, value)
}

// AttribFlag writes a name="name" attribute if flag is true,
// else nothing will be written.
func (self *XMLWriter) AttribFlag(name string, flag bool) *XMLWriter {
	if flag {
		self.Attrib(name, name)
	}
	return self
}

func (self *XMLWriter) Content(s string) *XMLWriter {
	self.Write([]byte(s))
	return self
}

func (self *XMLWriter) EscapeContent(s string) *XMLWriter {
	self.Write([]byte(html.EscapeString(s)))
	return self
}

func (self *XMLWriter) Printf(format string, args ...interface{}) *XMLWriter {
	fmt.Fprintf(self, format, args...)
	return self
}

func (self *XMLWriter) PrintfEscape(format string, args ...interface{}) *XMLWriter {
	return self.EscapeContent(fmt.Sprintf(format, args...))
}

// implements io.Writer
func (self *XMLWriter) Write(p []byte) (n int, err error) {
	self.finishOpenTag()
	return self.writer.Write(p)
}

func (self *XMLWriter) CloseTag() *XMLWriter {
	// this kind of sucks
	// if we can haz append() why not pop()?
	top := len(self.tagStack) - 1
	tag := self.tagStack[top]
	self.tagStack = self.tagStack[:top]

	if self.inOpenTag {
		self.writer.Write([]byte("/>"))
		self.inOpenTag = false
	} else {
		self.writer.Write([]byte("</"))
		self.writer.Write([]byte(tag))
		self.writer.Write([]byte{'>'})
	}

	return self
}

// Creates an explicit close tag, even if there is no content
func (self *XMLWriter) CloseTagAlways() *XMLWriter {
	self.finishOpenTag()
	return self.CloseTag()
}

func (self *XMLWriter) finishOpenTag() {
	if self.inOpenTag {
		self.writer.Write([]byte{'>'})
		self.inOpenTag = false
	}
}

func (self *XMLWriter) Reset() {
	if self.tagStack != nil {
		self.tagStack = self.tagStack[0:0]
	}
	self.inOpenTag = false
}
