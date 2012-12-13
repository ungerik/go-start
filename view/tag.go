package view

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

func (self *Tag) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag(self.Tag)
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	for key, value := range self.Attribs {
		ctx.Response.XML.Attrib(key, value)
	}
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	if self.ExtraClose {
		ctx.Response.XML.CloseTagAlways()
	} else {
		ctx.Response.XML.CloseTag()
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
