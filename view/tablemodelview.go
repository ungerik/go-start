package view

<<<<<<< HEAD
=======
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

>>>>>>> master
///////////////////////////////////////////////////////////////////////////////
// TableModelView

type TableModelView struct {
	ViewBase
	Class             string
	Caption           string
	GetModelIterator  GetModelIteratorFunc
	GetHeaderRowViews func(response *Response) (views Views, err error)
	GetRowViews       func(row int, rowModel interface{}, response *Response) (views Views, err error)
	table             Table
}

func (self *TableModelView) IterateChildren(callback IterateChildrenCallback) {
	callback(self, &self.table)
}

func (self *TableModelView) Render(response *Response) (err error) {
	self.table.Class = self.Class
	self.table.Caption = self.Caption

	tableModel := ViewsTableModel{}

	self.table.HeaderRow = false
	if self.GetHeaderRowViews != nil {
		views, err := self.GetHeaderRowViews(response)
		if err != nil {
			return err
		}
		if views != nil {
			tableModel = append(tableModel, views)
			self.table.HeaderRow = true
		}
	}

	rowNr := 0
	iter := self.GetModelIterator(response)
	for rowModel := iter.Next(); rowModel != nil; rowModel = iter.Next() {
		views, err := self.GetRowViews(rowNr, rowModel, response)
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

	return self.table.Render(response)
}
