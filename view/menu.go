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

func (self *Menu) Render(response *Response) (err error) {
	response.XML.OpenTag("ul")
	response.XML.AttribIfNotDefault("id", self.id)
	response.XML.AttribIfNotDefault("class", self.Class)

	requestURL := response.Request.URLString()

	// Find active item
	activeIndex := -1

	if self.ActiveItemClass != "" {
		// First try exact URL match
		for i := range self.Items {
			url := self.Items[i].URL(response)
			if url == requestURL {
				activeIndex = i
				break
			}
		}

		// If no exact URL match is found, search for sub pages
		if activeIndex == -1 {
			for i := range self.Items {
				url := self.Items[i].URL(response)
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
		url := linkModel.URL(response)

		// use i instead of index
		if i == activeIndex {
			itemClass += " " + self.ActiveItemClass
		}

		response.XML.OpenTag("li")
		if self.id != "" {
			response.XML.Attrib("id", self.id, "_", index)
		}
		response.XML.AttribIfNotDefault("class", itemClass)

		if i > 0 && self.BetweenItems != "" {
			response.XML.Content(self.BetweenItems)
		}

		response.XML.OpenTag("a")
		response.XML.Attrib("href", url)
		response.XML.AttribIfNotDefault("title", linkModel.LinkTitle(response))
		response.XML.AttribIfNotDefault("rel", linkModel.LinkRel(response))
		content := linkModel.LinkContent(response)
		if content != nil {
			err = content.Render(response)
			if err != nil {
				return err
			}
		}
		response.XML.ForceCloseTag() // a

		response.XML.ForceCloseTag() // li
	}

	response.XML.ForceCloseTag() // ul
	return nil
}
