package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// ShortTag

// ShortTag represents an arbitrary HTML element. It has a smaller footprint than Tag.
type ShortTag struct {
	ViewBase
	Tag     string
	Class   string
	Attribs map[string]string
	Content View
}

func (self *ShortTag) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *ShortTag) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag(self.Tag).AttribIfNotDefault("class", self.Class)
	for key, value := range self.Attribs {
		writer.Attrib(key, value)
	}
	if self.Content != nil {
		err = self.Content.Render(context, writer)
	}
	writer.ForceCloseTag()
	return err
}
