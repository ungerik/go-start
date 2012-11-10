package views

import (
	"github.com/ungerik/go-start/user"
	. "github.com/ungerik/go-start/view"
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

func Footer() View {
	return DIV("footer")
}

func TitleBar(title string) View {
	return &Div{Class: "title-bar", Content: Escape(title)}
}

func TitleBarRight(title string) View {
	return &Div{Class: "title-bar right", Content: Escape(title)}
}

func NewPublicPage(title string, main View) *Page {
	return &Page{
		Title: Escape(title),
		Scripts: Renderers{
			JQuery,
		},
		Content: Views{
			&Div{
				Class: "header",
				Content: Views{
					HeaderTopNav(),
					&Div{
						Class: "menu-area",
						Content: &Div{
							Class: "center",
							Content: Views{
								DIV("logo-container", IMG("/images/gopher.png")),
								HeaderMenu(),
							},
						},
					},
				},
			},
			&Div{
				Class: "content",
				Content: Views{
					&Div{
						Class: "center",
						Content: Views{
							&Div{
								Class:   "main",
								Content: main,
							},
							DivClearBoth(),
						},
					},
				},
			},
			Footer(),
		},
	}
}
