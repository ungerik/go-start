package view

import (
	"bytes"
	// "github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
)

///////////////////////////////////////////////////////////////////////////////
// Response

// not used atm.
type Response struct {
	buffer     bytes.Buffer
	webContext *web.Context
	// View that responds to the HTTP request
	RespondingView View
	// Custom response wide data that can be set by the application
	Data interface{}
}

func (self *Response) Write(p []byte) (n int, err error) {
	return self.buffer.Write(p)
}

func (self *Response) WriteByte(c byte) error {
	return self.buffer.WriteByte(c)
}

func (self *Response) WriteRune(r rune) (n int, err error) {
	return self.buffer.WriteRune(r)
}

func (self *Response) WriteString(s string) (n int, err error) {
	return self.buffer.WriteString(s)
}

func (self *Response) String() string {
	return self.buffer.String()
}

func (self *Response) SetSecureCookie(name string, val string, age int64, path string) {
	self.webContext.SetSecureCookie(name, val, age, path)
}

func (self *Response) Abort(status int, body string) {
	self.webContext.Abort(status, body)
}

func (self *Response) RedirectPermanently301(url string) {
	self.webContext.Redirect(301, url)
}

func (self *Response) RedirectTemporary302(url string) {
	self.webContext.Redirect(302, url)
}

func (self *Response) NotModified304(url string) {
	self.webContext.NotModified()
}

func (self *Response) NotFound404(url string) {
	self.Abort(404, "Not found")
}

func (self *Response) AuthorizationRequired401() {
	self.Abort(401, "Authorization Required")
}

func (self *Response) SetHeader(header string, value string, unique bool) {
	self.webContext.SetHeader(header, value, unique)
}
