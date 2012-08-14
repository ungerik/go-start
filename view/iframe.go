package view

///////////////////////////////////////////////////////////////////////////////
// Iframe

type Iframe struct {
	ViewBaseWithId
	Class        string
	Width        int
	Height       int
	Border       int
	Scrolling    bool
	MarginWidth  int
	MarginHeight int
	Seamless     bool
	URL          string
}

func (self *Iframe) Render(response *Response) (err error) {
	response.XML.OpenTag("iframe")
	response.XML.AttribIfNotDefault("id", self.id)
	response.XML.AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("width", self.Width).Attrib("height", self.Height)
	response.XML.Attrib("frameborder", self.Border)
	response.XML.Attrib("marginwidth", self.MarginWidth).Attrib("marginheight", self.MarginHeight)
	if self.Scrolling {
		response.XML.Attrib("scrolling", "yes")
	} else {
		response.XML.Attrib("scrolling", "no")
	}
	if self.Seamless {
		response.XML.Attrib("seamless", "seamless")
	}
	response.XML.Attrib("src", self.URL)
	response.XML.ForceCloseTag()
	return nil
}
