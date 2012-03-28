package view

import (
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
)

///////////////////////////////////////////////////////////////////////////////
// Response

// not used atm.
type Response struct {
	utils.XMLBuffer
	webContext *web.Context
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
	self.webContext.NotFound("Not found")
}
