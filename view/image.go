package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Image

type Image struct {
	ViewBaseWithId
	Class       string
	URL         string
	Width       int
	Height      int
	Description string
}

func (self *Image) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("img").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("src", self.URL)
	writer.AttribIfNotDefault("width", self.Width)
	writer.AttribIfNotDefault("height", self.Height)
	writer.AttribIfNotDefault("alt", self.Description)
	writer.CloseTag()
	return nil
}

//func (self *Image) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
