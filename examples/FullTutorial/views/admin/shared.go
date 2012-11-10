package admin

import (
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
	. "github.com/ungerik/go-start/view"
)

func NewAdminPage(title string, main View) *Page {
	return &Page{
		Title: Escape(title),
		// CSS:     IndirectURL(&Admin_CSS),
		Scripts: Renderers{
			JQuery,
		},
		Content: Views{
			&Div{
				Class: "header",
				Content: Views{
					&Link{
						Class: "title",
						Model: &PageLink{
							Page: &Admin,
							Content: &Tag{
								Tag: "h1",
								Content: Views{
									&Image{Class: "logo", Src: "/images/gopher.png"},
									HTML("Admin Panel"),
								},
							},
						},
					},
					HeaderUserNav(),
					&Div{
						Class: "menu-frame",
						Content: Views{
							&Menu{
								Class:           "menu",
								ItemClass:       "menu-item",
								ActiveItemClass: "active",
								BetweenItems:    " &nbsp;/&nbsp; ",
								Items: []LinkModel{
									NewPageLink(&Admin, "Dashboard"),
									NewPageLink(&Admin_Users, "Users"),
									NewPageLink(&Admin_Images, "Images"),
								},
							},
							DivClearBoth(),
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
