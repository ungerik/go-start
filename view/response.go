package view

import (
	"bytes"
	"fmt"
	"github.com/ungerik/web.go"
)

func newResponse(webContext *web.Context, respondingView View) *Response {
	return &Response{
		webContext:     webContext,
		RespondingView: respondingView,
	}
}

type Response struct {
	buffer     bytes.Buffer
	webContext *web.Context
	// View that responds to the HTTP request
	RespondingView View
	// Custom response wide data that can be set by the application
	Data interface{}
}

// New creates a clone of the response with an empty buffer.
// Used to render preliminary text.
func (self *Response) New() *Response {
	return &Response{
		webContext:     self.webContext,
		RespondingView: self.RespondingView,
		Data:           self.Data,
	}
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

func (self *Response) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(&self.buffer, format, args...)
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

func (self *Response) NotModified304() {
	self.webContext.NotModified()
}

func (self *Response) Forbidden403(message string) {
	self.Abort(403, message)
}

func (self *Response) NotFound404(message string) {
	self.Abort(404, message)
}

func (self *Response) AuthorizationRequired401() {
	self.Abort(401, "Authorization Required")
}

func (self *Response) SetHeader(header string, value string, unique bool) {
	self.webContext.SetHeader(header, value, unique)
}

func (self *Response) ContentType(ext string) {
	self.webContext.ContentType(ext)
}
