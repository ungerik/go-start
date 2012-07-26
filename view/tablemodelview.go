package view

import (
	"github.com/ungerik/go-start/utils"
)

func TableHeaderRow(views ...View) func(context *Context) (Views, error) {
	return func(context *Context) (Views, error) {
		return Views(views), nil
	}
}

func TableHeaderRowEscape(s ...string) func(context *Context) (Views, error) {
	views := make(Views, len(s))
	for i := range s {
		views[i] = Escape(s[i])
	}
	return func(context *Context) (Views, error) {
		return views, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// TableModelIteratorView

type TableModelIteratorView struct {
	ViewBase
	Class            string
	Caption          string
	GetModelIterator GetModelIteratorFunc
	GetHeaderRow     func(context *Context) (views Views, err error)
	GetRow           func(row int, rowModel interface{}, context *Context) (views Views, err error)
	table            Table
}

func (self *TableModelIteratorView) IterateChildren(callback IterateChildrenCallback) {
	callback(self, &self.table)
}

func (self *TableModelIteratorView) Render(context *Context, writer *utils.XMLWriter) (err error) {
	self.table.Class = self.Class
	self.table.Caption = self.Caption

	tableModel := ViewsTableModel{}

	self.table.HeaderRow = false
	if self.GetHeaderRow != nil {
		views, err := self.GetHeaderRow(context)
		if err != nil {
			return err
		}
		if views != nil {
			tableModel = append(tableModel, views)
			self.table.HeaderRow = true
		}
	}

	rowNr := 0
	iter := self.GetModelIterator(context)
	for rowModel := iter.Next(); rowModel != nil; rowModel = iter.Next() {
		views, err := self.GetRow(rowNr, rowModel, context)
		if err != nil {
			return err
		}
		if views != nil {
			tableModel = append(tableModel, views)
			rowNr++
		}
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	self.table.Model = tableModel

	return self.table.Render(context, writer)
}
