package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Table

type Table struct {
	ViewBaseWithId
	Model     TableModel
	Class     string
	Caption   string
	HeaderRow bool
}

func (self *Table) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("table").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)

	if self.Caption != "" {
		writer.OpenTag("caption").EscapeContent(self.Caption).CloseTag()
	}

	if self.Model != nil {
		rows := self.Model.Rows()
		columns := self.Model.Columns()

		for row := 0; row < rows; row++ {
			writer.OpenTag("tr")
			if row&1 == 0 {
				writer.Attrib("class", "row", row, " even")
			} else {
				writer.Attrib("class", "row", row, " odd")
			}

			for col := 0; col < columns; col++ {
				if self.HeaderRow && row == 0 {
					writer.OpenTag("th")
				} else {
					writer.OpenTag("td")
				}
				if col&1 == 0 {
					writer.Attrib("class", "col", col, " even")
				} else {
					writer.Attrib("class", "col", col, " odd")
				}
				view, err := self.Model.CellView(row, col, context)
				if view != nil && err == nil {
					view.Init(view)
					err = view.Render(context, writer)
				}
				if err != nil {
					return err
				}
				writer.ExtraCloseTag() // td/th
			}

			writer.ExtraCloseTag() // tr
		}
	}

	writer.ExtraCloseTag() // table
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// TableModel

type TableModel interface {
	Rows() int
	Columns() int
	CellView(row int, column int, context *Context) (view View, err error)
}

///////////////////////////////////////////////////////////////////////////////
// ViewsTableModel

type ViewsTableModel []Views

func (self ViewsTableModel) Rows() int {
	return len(self)
}

func (self ViewsTableModel) Columns() int {
	if len(self) == 0 {
		return 0
	}
	return len(self[0])
}

func (self ViewsTableModel) CellView(row int, column int, context *Context) (view View, err error) {
	return self[row][column], nil
}

///////////////////////////////////////////////////////////////////////////////
// StringsTableModel

type HTMLStringsTableModel [][]string

func (self HTMLStringsTableModel) Rows() int {
	return len(self)
}

func (self HTMLStringsTableModel) Columns() int {
	if len(self) == 0 {
		return 0
	}
	return len(self[0])
}

func (self HTMLStringsTableModel) CellView(row int, column int, context *Context) (view View, err error) {
	return HTML(self[row][column]), nil
}

///////////////////////////////////////////////////////////////////////////////
// StringsTableModel

type EscapeStringsTableModel [][]string

func (self EscapeStringsTableModel) Rows() int {
	return len(self)
}

func (self EscapeStringsTableModel) Columns() int {
	if len(self) == 0 {
		return 0
	}
	return len(self[0])
}

func (self EscapeStringsTableModel) CellView(row int, column int, context *Context) (view View, err error) {
	return Escape(self[row][column]), nil
}
