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
			DIV("header",
				&Link{
					Class: "title",
					Model: &PageLink{
						Page:    &Admin,
						Content: H1(&Image{Class: "logo", Src: "/images/gopher.png"}, HTML("Admin Panel")),
					},
				},
				HeaderUserNav(),
				DIV("menu-frame",
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
				),
			),
			DIV("content",
				DIV("center",
					DIV("main", main),
					DivClearBoth(),
				),
			),
			Footer(),
		},
	}
}
