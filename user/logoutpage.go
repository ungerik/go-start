package user

import (
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/go-start/view"
)

// If redirectURL is nil, the redirect will go to "/"
func NewLogoutView(redirectURL view.URL) view.ViewWithURL {
	return view.NewViewURLWrapper(view.Renderer(
		func(context *view.Context, writer *utils.XMLWriter) (err error) {
			Logout(context)
			url := "/"
			if redirectURL != nil {
				url = redirectURL.URL(context)
			}
			return view.Redirect(url)
		},
	))
}
