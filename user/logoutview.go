package user

import (
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/go-start/view"
)

// If redirect is nil, the redirect will go to "/"
func LogoutView(redirect view.URL) view.View {
	return view.RenderView(
		func(context *view.Context, writer *utils.XMLWriter) (err error) {
			Logout(context)
			if redirect != nil {
				return view.Redirect(redirect.URL(context))
			}
			return view.Redirect("/")
		},
	)
}
