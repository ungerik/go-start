package view

import "github.com/ungerik/go-start/utils"

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

func (self *Menu) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("ul").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)

	requestURL := context.RequestURL()

	// Find active item
	activeIndex := -1

	if self.ActiveItemClass != "" {
		// First try exact URL match
		for i := range self.Items {
			url := self.Items[i].URL(context)
			if url == requestURL {
				activeIndex = i
				break
			}
		}

		// If no exact URL match is found, search for sub pages
		if activeIndex == -1 {
			for i := range self.Items {
				url := self.Items[i].URL(context)
				if utils.StringStartsWith(requestURL, url) {
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
		url := linkModel.URL(context)

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
		writer.AttribIfNotDefault("title", linkModel.LinkTitle(context))
		writer.AttribIfNotDefault("rel", linkModel.LinkRel(context))
		content := linkModel.LinkContent(context)
		if content != nil {
			err = content.Render(context, writer)
			if err != nil {
				return err
			}
		}
		writer.ExtraCloseTag() // a

		writer.ExtraCloseTag() // li
	}

	writer.ExtraCloseTag() // ul
	return nil
}
