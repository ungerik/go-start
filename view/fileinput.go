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

func (self *FileInput) Render(response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("type", "file").Attrib("name", self.Name)
	if self.Disabled {
		writer.Attrib("disabled", "disabled")
	}
	writer.CloseTag()
	return err
}
