package views

import (
	. "github.com/ungerik/go-start/view"
)

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
