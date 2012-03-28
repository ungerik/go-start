package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Span

// Span represents a HTML span element.
type Span struct {
	ViewBaseWithId
	Class   string
	Content View
}

func (self *Span) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Span) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("span").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	if self.Content != nil {
		err = self.Content.Render(context, writer)
	}
	writer.ExtraCloseTag()
	return err
}
