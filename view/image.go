package view

///////////////////////////////////////////////////////////////////////////////
// Image

type Image struct {
	ViewBaseWithId
	Class       string
	URL         URL    // If URL is set, then Src will be ignored
	Src         string // String URL of the image, used when URL is nil
	Width       int
	Height      int
	Description string
}

func (self *Image) Render(response *Response) (err error) {
	response.XML.OpenTag("img").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	src := self.Src
	if self.URL != nil {
		src = self.URL.URL(response)
	}
	response.XML.Attrib("src", src)
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
