package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// List

// TODO definition list
type List struct {
	ViewBaseWithId
	Model       ListModel
	Ordered     bool
	OrderOffset uint
	Class       string
}

func (self *List) Render(response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
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
			view, err := self.Model.ItemView(i, response)
			if view != nil && err == nil {
				view.Init(view)
				err = view.Render(response)
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
	ItemView(index int, response *Response) (view View, err error)
}

///////////////////////////////////////////////////////////////////////////////
// ViewsListModel

type ViewsListModel []View

func (self ViewsListModel) NumItems() int {
	return len(self)
}

func (self ViewsListModel) ItemView(index int, response *Response) (view View, err error) {
	return self[index], nil
}

///////////////////////////////////////////////////////////////////////////////
// MultiViewsListModel

type MultiViewsListModel []Views

func (self MultiViewsListModel) NumItems() int {
	return len(self)
}

func (self MultiViewsListModel) ItemView(index int, response *Response) (view View, err error) {
	return self[index], nil
}

///////////////////////////////////////////////////////////////////////////////
// HTMLStringsListModel

type HTMLStringsListModel []string

func (self HTMLStringsListModel) NumItems() int {
	return len(self)
}

func (self HTMLStringsListModel) ItemView(index int, response *Response) (view View, err error) {
	return HTML(self[index]), nil
}

///////////////////////////////////////////////////////////////////////////////
// EscapeStringsListModel

type EscapeStringsListModel []string

func (self EscapeStringsListModel) NumItems() int {
	return len(self)
}

func (self EscapeStringsListModel) ItemView(index int, response *Response) (view View, err error) {
	return Escape(self[index]), nil
}
