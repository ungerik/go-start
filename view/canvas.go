package view

import "github.com/ungerik/go-start/utils"

type Canvas struct {
	ViewBaseWithId
	Class  string
	Width  int
	Height int
}

func (self *Canvas) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("label").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("width", self.Width).Attrib("height", self.Height)
	writer.ExtraCloseTag()
	return err
}
