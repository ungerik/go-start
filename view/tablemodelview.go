package view

import

//	"github.com/ungerik/go-start/model"
"github.com/ungerik/go-start/utils"

//	"github.com/ungerik/go-start/debug"

///////////////////////////////////////////////////////////////////////////////
// TableModelView

type TableModelView struct {
	ViewBase
	Class             string
	Caption           string
	GetModelIterator  GetModelIteratorFunc
	GetHeaderRowViews func(context *Context) (views Views, err error)
	GetRowViews       func(row int, rowModel interface{}, context *Context) (views Views, err error)
	table             Table
}

func (self *TableModelView) IterateChildren(callback IterateChildrenCallback) {
	callback(self, &self.table)
}

func (self *TableModelView) Render(context *Context, writer *utils.XMLWriter) (err error) {
	self.table.Class = self.Class
	self.table.Caption = self.Caption

	tableModel := ViewsTableModel{}

	self.table.HeaderRow = false
	if self.GetHeaderRowViews != nil {
		views, err := self.GetHeaderRowViews(context)
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
		views, err := self.GetRowViews(rowNr, rowModel, context)
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
