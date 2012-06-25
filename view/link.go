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

func (self *Link) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("a").AttribIfNotDefault("id", self.id)
	writer.AttribIfNotDefault("class", self.Class)
	if self.NewWindow {
		writer.Attrib("target", "_blank")
	}
	if self.Model != nil {
		writer.Attrib("href", self.Model.URL(context))
		writer.AttribIfNotDefault("title", self.Model.LinkTitle(context))
		writer.AttribIfNotDefault("rel", self.Model.LinkRel(context))
		content := self.Model.LinkContent(context)
		if content != nil {
			err = content.Render(context, writer)
		}
	}
	writer.ForceCloseTag() // a
	return err
}
