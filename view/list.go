package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// List

// TODO definition list
type List struct {
	ViewBaseWithIdAndDynamicChildren
	Model       ListModel
	Ordered     bool
	OrderOffset uint
	Class       string
}

func (self *List) Render(context *Context, writer *utils.XMLWriter) (err error) {
	self.RemoveChildren()

	if self.Ordered {
		writer.OpenTag("ol").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		writer.Attrib("start", self.OrderOffset+1)
	} else {
		writer.OpenTag("ul").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	}

	if self.Model != nil {
		numItems := self.Model.NumItems()
		for i := 0; i < numItems; i++ {
			writer.OpenTag("li").Attrib("id", self.id, "_", i)
			view, err := self.Model.ItemView(i, context)
			if view != nil && err == nil {
				self.AddAndInitChild(view)
				err = view.Render(context, writer)
			}
			if err != nil {
				return err
			}
			writer.ExtraCloseTag() // li
		}
	}

	writer.ExtraCloseTag() // ol/ul
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ListModel

type ListModel interface {
	NumItems() int
	ItemView(index int, context *Context) (view View, err error)
}

///////////////////////////////////////////////////////////////////////////////
// ViewsListModel

type ViewsListModel []View

func (self ViewsListModel) NumItems() int {
	return len(self)
}

func (self ViewsListModel) ItemView(index int, context *Context) (view View, err error) {
	return self[index], nil
}

///////////////////////////////////////////////////////////////////////////////
// MultiViewsListModel

type MultiViewsListModel []Views

func (self MultiViewsListModel) NumItems() int {
	return len(self)
}

func (self MultiViewsListModel) ItemView(index int, context *Context) (view View, err error) {
	return self[index], nil
}

///////////////////////////////////////////////////////////////////////////////
// HTMLStringsListModel

type HTMLStringsListModel []string

func (self HTMLStringsListModel) NumItems() int {
	return len(self)
}

func (self HTMLStringsListModel) ItemView(index int, context *Context) (view View, err error) {
	return HTML(self[index]), nil
}

///////////////////////////////////////////////////////////////////////////////
// EscapeStringsListModel

type EscapeStringsListModel []string

func (self EscapeStringsListModel) NumItems() int {
	return len(self)
}

func (self EscapeStringsListModel) ItemView(index int, context *Context) (view View, err error) {
	return Escape(self[index]), nil
}
