package subpage

import (
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/ViewPaths/views"
	"github.com/ungerik/go-start/examples/ViewPaths/views/root"
)

func init() {
	views.Admin_User0 = &Page{
		Title: RenderView(
			func(ctx *Context) error {
				// The username is in ctx.URLArgs[0]
				ctx.Response.WriteString("Manage " + ctx.URLArgs[0])
				return nil
			},
		),
		Content: Views{
			DynamicViewBindURLArgs(
				// The URL argument 0 can also be bound dynamically
				// to a function argument:
				func(ctx *Context, username string) (View, error) {
					return H1("Manage user ", username), nil
				},
			),
			root.Navigation(),
			DynamicViewBindURLArgs(
				func(ctx *Context, username string) (View, error) {
					return Views{
						H4(Printf("This view uses the URL argument '%s':", username)),
						HTML(views.Admin_User0.URL(ctx)),
					}, nil
				},
			),
		},
	}
}
