package utils

import (
	"bytes"
)

///////////////////////////////////////////////////////////////////////////////
// XMLBuffer

func NewXMLBuffer() *XMLBuffer {
	var xmlBuffer XMLBuffer
	xmlBuffer.writer = &xmlBuffer.buffer
	return &xmlBuffer
}

type XMLBuffer struct {
	XMLWriter
	buffer bytes.Buffer
}

func (self *XMLBuffer) String() string {
	return self.buffer.String()
}

func (self *XMLBuffer) Bytes() []byte {
	return self.buffer.Bytes()
}
