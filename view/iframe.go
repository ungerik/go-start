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

func (self *Iframe) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("iframe")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.Attrib("width", self.Width).Attrib("height", self.Height)
	ctx.Response.XML.Attrib("frameborder", self.Border)
	ctx.Response.XML.Attrib("marginwidth", self.MarginWidth).Attrib("marginheight", self.MarginHeight)
	if self.Scrolling {
		ctx.Response.XML.Attrib("scrolling", "yes")
	} else {
		ctx.Response.XML.Attrib("scrolling", "no")
	}
	if self.Seamless {
		ctx.Response.XML.Attrib("seamless", "seamless")
	}
	ctx.Response.XML.Attrib("src", self.URL)
	ctx.Response.XML.CloseTagAlways()
	return nil
}
