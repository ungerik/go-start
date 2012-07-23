package view

import (
	"github.com/ungerik/go-start/utils"
)

type FileInput struct {
	ViewBaseWithId
	Class    string
	Name     string
	Disabled bool
}

func (self *FileInput) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("type", "file").Attrib("name", self.Name)
	if self.Disabled {
		writer.Attrib("disabled", "disabled")
	}
	writer.CloseTag()
	return err
}
