package view

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

func (self *Image) Render(response *Response) (err error) {
	response.XML.OpenTag("img").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("src", self.URL)
	response.XML.AttribIfNotDefault("width", self.Width)
	response.XML.AttribIfNotDefault("height", self.Height)
	response.XML.AttribIfNotDefault("alt", self.Description)
	response.XML.CloseTag()
	return nil
}

//func (self *Image) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
