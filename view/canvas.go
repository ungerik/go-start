package view

import "github.com/ungerik/go-start/utils"

type Canvas struct {
	ViewBaseWithId
	Class  string
	Width  int
	Height int
}

func (self *Canvas) Render(request *Request, session *Session, response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("label").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("width", self.Width).Attrib("height", self.Height)
	writer.ExtraCloseTag()
	return err
}
