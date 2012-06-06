package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Tag

// Tag represents an arbitrary HTML element.
type Tag struct {
	ViewBaseWithId
	Tag        string
	Content    View
	Class      string
	Attribs    map[string]string
	ExtraClose bool
}

func (self *Tag) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Tag) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag(self.Tag).Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	for key, value := range self.Attribs {
		writer.Attrib(key, value)
	}
	if self.Content != nil {
		err = self.Content.Render(context, writer)
	}
	if self.ExtraClose {
		writer.ForceCloseTag()
	} else {
		writer.CloseTag()
	}
	return err
}

//func (self *Tag) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
//
//func (self *Tag) SetContent(content View) {
//	self.Content = content
//	ViewChanged(self)
//}
//
//func (self *Tag) SetTag(tag string) {
//	self.Tag = tag
//	ViewChanged(self)
//}
