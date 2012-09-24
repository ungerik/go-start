package root

import (
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/ViewPaths/views"
)

func init() {
	views.GetXML = NewViewURLWrapper(RenderView(
		func(ctx *Context) error {
			ctx.Response.SetContentTypeByExt(".xml")
			// ctx.Response.XML is convenience writer for XML and HTML
			ctx.Response.XML.WriteXMLDeclaration()
			ctx.Response.XML.OpenTag("xml")
			ctx.Response.XML.OpenTag("data").Content("Hello World!").CloseTag()
			ctx.Response.XML.CloseTag() // xml
			return nil
		},
	))
}
