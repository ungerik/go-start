package view

import ()

type FileInput struct {
	ViewBaseWithId
	Class    string
	Name     string
	Disabled bool
}

func (self *FileInput) Render(response *Response) (err error) {
	response.XML.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("type", "file").Attrib("name", self.Name)
	if self.Disabled {
		response.XML.Attrib("disabled", "disabled")
	}
	response.XML.CloseTag()
	return err
}
