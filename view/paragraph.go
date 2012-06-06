package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Paragraph

type Paragraph struct {
	ViewBaseWithId
	Class   string
	Content View
}

func (self *Paragraph) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Paragraph) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("p").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	if self.Content != nil {
		err = self.Content.Render(context, writer)
	}
	writer.ForceCloseTag()
	return err
}
