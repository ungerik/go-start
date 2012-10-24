package view

func TableHeaderRow(views ...View) func(ctx *Context) (Views, error) {
	return func(ctx *Context) (Views, error) {
		return Views(views), nil
	}
}

func TableHeaderRowEscape(s ...string) func(ctx *Context) (Views, error) {
	views := make(Views, len(s))
	for i := range s {
		views[i] = Escape(s[i])
	}
	return func(ctx *Context) (Views, error) {
		return views, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// ModelIteratorTableView

type ModelIteratorTableView struct {
	ViewBase
	Class             string
	Caption           string
	GetModelIterator  GetModelIteratorFunc
	GetHeaderRowViews func(ctx *Context) (views Views, err error)
	// GetRowModel returns the model used GetRowViews.
	// Only return the same object if you know exactly what you are doing.
	// To be on the safe side, return a new model instance at every call.
	GetRowModel func(ctx *Context) (model interface{}, err error)
	GetRowViews func(row int, rowModel interface{}, ctx *Context) (views Views, err error)
	table       Table
}

func (self *ModelIteratorTableView) IterateChildren(callback IterateChildrenCallback) {
	callback(self, &self.table)
}

func (self *ModelIteratorTableView) Render(ctx *Context) (err error) {
	self.table.Class = self.Class
	self.table.Caption = self.Caption

	tableModel := ViewsTableModel{}

	self.table.HeaderRow = false
	if self.GetHeaderRowViews != nil {
		views, err := self.GetHeaderRowViews(ctx)
		if err != nil {
			return err
		}
		if views != nil {
			tableModel = append(tableModel, views)
			self.table.HeaderRow = true
		}
	}

	rowNr := 0
	iter := self.GetModelIterator(ctx)
	rowModel, err := self.GetRowModel(ctx)
	if err != nil {
		return err
	}
	for iter.Next(rowModel) {
		views, err := self.GetRowViews(rowNr, rowModel, ctx)
		if err != nil {
			return err
		}
		if views != nil {
			tableModel = append(tableModel, views)
			rowNr++
		}
		rowModel, err = self.GetRowModel(ctx)
		if err != nil {
			return err
		}
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	self.table.Model = tableModel

	return self.table.Render(ctx)
}
