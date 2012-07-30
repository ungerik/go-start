package view

import (
	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// Link

type Link struct {
	ViewBaseWithId
	Class     string
	Model     LinkModel
	NewWindow bool
}

func (self *Link) Render(response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("a").AttribIfNotDefault("id", self.id)
	writer.AttribIfNotDefault("class", self.Class)
	if self.NewWindow {
		writer.Attrib("target", "_blank")
	}
	if self.Model != nil {
		writer.Attrib("href", self.Model.URL(response.Request.URLArgs...))
		writer.AttribIfNotDefault("title", self.Model.LinkTitle(response))
		writer.AttribIfNotDefault("rel", self.Model.LinkRel(response))
		content := self.Model.LinkContent(response)
		if content != nil {
			err = content.Render(response)
		}
	}
	writer.ForceCloseTag() // a
	return err
}
