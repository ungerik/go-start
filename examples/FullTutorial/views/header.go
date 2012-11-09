package views

import (
	// "github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/user"
	. "github.com/ungerik/go-start/view"

	// "github.com/ungerik/go-start/examples/FullTutorial/models"
)

func HeaderTopNav() View {
	return &Div{
		Class: "top-nav",
		Content: &Div{
			Class: "center",
			Content: Views{
				HeaderUserNav(),
				DynamicView(
					func(ctx *Context) (View, error) {
						if ctx.Request.RequestURI == "/" {
							return nil, nil
						}
						return A("/", HTML("&larr; Back to the homepage")), nil
					},
				),
			},
		},
	}
}

func HeaderMenu() *Menu {
	return &Menu{
		Class:           "menu",
		ItemClass:       "menu-item",
		ActiveItemClass: "active",
		Items: []LinkModel{
			NewLinkModel(&Homepage, "Home"),
		},
	}
}

func HeaderUserNav() View {
	return DIV("login-nav",
		user.Nav(
			A(&LoginSignup, "Login / Sign up"),
			nil,
			A(&Logout, "Logout"),
			A(&Profile, "My profile"),
			HTML("&nbsp; | &nbsp;"),
		),
	)
}
