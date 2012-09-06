package view

///////////////////////////////////////////////////////////////////////////////
// Image

type Image struct {
	ViewBaseWithId
	Class  string
	URL    URL    // If URL is set, then Src will be ignored
	Src    string // String URL of the image, used when URL is nil
	Width  int
	Height int
	Title  string
}

func (self *Image) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("img")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	src := self.Src
	if self.URL != nil {
		src = self.URL.URL(ctx)
	}
	ctx.Response.XML.Attrib("src", src)
	ctx.Response.XML.AttribIfNotDefault("width", self.Width)
	ctx.Response.XML.AttribIfNotDefault("height", self.Height)
	ctx.Response.XML.AttribIfNotDefault("alt", self.Title)
	ctx.Response.XML.CloseTag()
	return nil
}
