package view

import (
	"strings"
	"github.com/ungerik/go-start/utils"
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
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("ul").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)

	requestURL := response.Request.URLString()

	// Find active item
	activeIndex := -1

	if self.ActiveItemClass != "" {
		// First try exact URL match
		for i := range self.Items {
			url := self.Items[i].URL(response.Request.PathArgs...)
>>>>>>> master
			if url == requestURL {
				activeIndex = i
				break
			}
		}

		// If no exact URL match is found, search for sub pages
		if activeIndex == -1 {
			for i := range self.Items {
<<<<<<< HEAD
				url := self.Items[i].URL(response)
				if utils.StringStartsWith(requestURL, url) {
=======
				url := self.Items[i].URL(context.PathArgs...)
				if strings.HasPrefix(requestURL, url) {
>>>>>>> master
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
		url := linkModel.URL(response.Request.Params...)

		// use i instead of index
		if i == activeIndex {
			itemClass += " " + self.ActiveItemClass
		}

		writer.OpenTag("li").Attrib("id", self.id, "_", index).AttribIfNotDefault("class", itemClass)

		if i > 0 && self.BetweenItems != "" {
			writer.Content(self.BetweenItems)
		}

		writer.OpenTag("a")
		writer.Attrib("href", url)
		writer.AttribIfNotDefault("title", linkModel.LinkTitle(response))
		writer.AttribIfNotDefault("rel", linkModel.LinkRel(response))
		content := linkModel.LinkContent(response)
		if content != nil {
			err = content.Render(response)
			if err != nil {
				return err
			}
		}
		writer.ForceCloseTag() // a

		writer.ForceCloseTag() // li
	}

	writer.ForceCloseTag() // ul
	return nil
}
