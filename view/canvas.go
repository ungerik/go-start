package view

type Canvas struct {
	ViewBaseWithId
	Class  string
	Width  int
	Height int
}

func (self *Canvas) Render(response *Response) (err error) {
	response.XML.OpenTag("label").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("width", self.Width).Attrib("height", self.Height)
	response.XML.ForceCloseTag()
	return err
}
