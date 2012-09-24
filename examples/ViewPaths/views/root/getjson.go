package root

import (
	. "github.com/ungerik/go-start/view"
	
	"github.com/ungerik/go-start/examples/ViewPaths/views"
)

func init() {
	views.GetJSON = NewViewURLWrapper(RenderView(
		func(ctx *Context) error {
			ctx.Response.SetContentTypeByExt(".json")
			ctx.Response.WriteString(`{"data": "Hello World!"}`)
			return nil
		},
	))
}
