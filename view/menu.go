package view

import (
	"strings"
	// "github.com/ungerik/go-start/debug"
)

///////////////////////////////////////////////////////////////////////////////
// Menu

type Menu struct {
	ViewBaseWithId
	Class           string
	ItemClass       string
	ActiveItemClass string
	BetweenItems    string
	Items           []LinkModel
	Reverse         bool
}

func (self *Menu) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("ul")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)

	requestURL := ctx.Request.URLString()

	// Find active item
	activeIndex := -1

	if self.ActiveItemClass != "" {
		// First try exact URL match
		for i := range self.Items {
			url := self.Items[i].URL(ctx)
			if url == requestURL {
				activeIndex = i
				break
			}
		}

		// If no exact URL match is found, search for sub pages
		if activeIndex == -1 {
			for i := range self.Items {
				url := self.Items[i].URL(ctx)
				if strings.HasPrefix(requestURL, url) {
					activeIndex = i
					// todo
					// not perfect, what if homepage matches first, but other matches better?
				}
			}
		}
	}

	for i := range self.Items {
		index := i
		if self.Reverse {
			index = len(self.Items) - 1 - i
		}
		itemClass := self.ItemClass
		linkModel := self.Items[index]
		url := linkModel.URL(ctx)

		// use i instead of index
		if i == activeIndex {
			itemClass += " " + self.ActiveItemClass
		}

		ctx.Response.XML.OpenTag("li")
		if self.id != "" {
			ctx.Response.XML.Attrib("id", self.id, "_", index)
		}
		ctx.Response.XML.AttribIfNotDefault("class", itemClass)

		if i > 0 && self.BetweenItems != "" {
			ctx.Response.XML.Content(self.BetweenItems)
		}

		ctx.Response.XML.OpenTag("a")
		ctx.Response.XML.Attrib("href", url)
		ctx.Response.XML.AttribIfNotDefault("title", linkModel.LinkTitle(ctx))
		ctx.Response.XML.AttribIfNotDefault("rel", linkModel.LinkRel(ctx))
		content := linkModel.LinkContent(ctx)
		if content != nil {
			err = content.Render(ctx)
			if err != nil {
				return err
			}
		}
		ctx.Response.XML.CloseTagAlways() // a

		ctx.Response.XML.CloseTagAlways() // li
	}

	ctx.Response.XML.CloseTagAlways() // ul
	return nil
}
