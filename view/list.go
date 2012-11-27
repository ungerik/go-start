package view

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

func (self *List) Render(ctx *Context) (err error) {
	if self.Ordered {
		ctx.Response.XML.OpenTag("ol")
		ctx.Response.XML.Attrib("start", self.OrderOffset+1)
	} else {
		ctx.Response.XML.OpenTag("ul")
	}
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)

	if self.Model != nil {
		numItems := self.Model.NumItems()
		for i := 0; i < numItems; i++ {
			ctx.Response.XML.OpenTag("li")
			if self.id != "" {
				ctx.Response.XML.Attrib("id", self.id, "_", i)
			}
			view, err := self.Model.ItemView(i, ctx)
			if view != nil && err == nil {
				view.Init(view)
				err = view.Render(ctx)
			}
			if err != nil {
				return err
			}
			ctx.Response.XML.CloseTagAlways() // li
		}
	}

	ctx.Response.XML.CloseTagAlways() // ol/ul
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ListModel

type ListModel interface {
	NumItems() int
	ItemView(index int, ctx *Context) (view View, err error)
}

///////////////////////////////////////////////////////////////////////////////
// ViewsListModel

type ViewsListModel []View

func (self ViewsListModel) NumItems() int {
	return len(self)
}

func (self ViewsListModel) ItemView(index int, ctx *Context) (view View, err error) {
	return self[index], nil
}

///////////////////////////////////////////////////////////////////////////////
// MultiViewsListModel

type MultiViewsListModel []Views

func (self MultiViewsListModel) NumItems() int {
	return len(self)
}

func (self MultiViewsListModel) ItemView(index int, ctx *Context) (view View, err error) {
	return self[index], nil
}

///////////////////////////////////////////////////////////////////////////////
// HTMLStringsListModel

type HTMLStringsListModel []string

func (self HTMLStringsListModel) NumItems() int {
	return len(self)
}

func (self HTMLStringsListModel) ItemView(index int, ctx *Context) (view View, err error) {
	return HTML(self[index]), nil
}

///////////////////////////////////////////////////////////////////////////////
// EscapeStringsListModel

type EscapeStringsListModel []string

func (self EscapeStringsListModel) NumItems() int {
	return len(self)
}

func (self EscapeStringsListModel) ItemView(index int, ctx *Context) (view View, err error) {
	return Escape(self[index]), nil
}
