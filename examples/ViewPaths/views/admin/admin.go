package admin

import (
	. "github.com/ungerik/go-start/view"
	
	"github.com/ungerik/go-start/examples/ViewPaths/views"
	"github.com/ungerik/go-start/examples/ViewPaths/views/root"
)

func init() {
	views.AdminAuth = NewBasicAuth("Username: admin, Password: admin", "admin", "admin")

	views.Admin = &Page{
		Title: Escape("Admin"),
		Content: Views{
			H1("Admin Page"),
			root.Navigation(),
			H3("Manage Users:"),
			UL(
				DynamicView(
					func(ctx *Context) (View, error) {
						url := views.Admin_User0.URL(ctx.ForURLArgs("ErikUnger"))
						return A(url, "ErikUnger"), nil
					},
				),
				DynamicView(
					func(ctx *Context) (View, error) {
						url := views.Admin_User0.URL(ctx.ForURLArgs("AlexTacho"))
						return A(url, "AlexTacho"), nil
					},
				),
			),
		},
	}
}
