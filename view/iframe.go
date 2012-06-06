package view

import "github.com/ungerik/go-start/utils"

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

func (self *Iframe) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("iframe").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("width", self.Width).Attrib("height", self.Height)
	writer.Attrib("frameborder", self.Border)
	writer.Attrib("marginwidth", self.MarginWidth).Attrib("marginheight", self.MarginHeight)
	if self.Scrolling {
		writer.Attrib("scrolling", "yes")
	} else {
		writer.Attrib("scrolling", "no")
	}
	if self.Seamless {
		writer.Attrib("seamless", "seamless")
	}
	writer.Attrib("src", self.URL)
	writer.ForceCloseTag()
	return nil
}
