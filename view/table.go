package view

///////////////////////////////////////////////////////////////////////////////
// Table

type Table struct {
	ViewBaseWithId
	Model     TableModel
	Class     string
	Caption   string
	HeaderRow bool
}

func (self *Table) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("table")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)

	if self.Caption != "" {
		ctx.Response.XML.OpenTag("caption").EscapeContent(self.Caption).CloseTag()
	}

	if self.Model != nil {
		rows := self.Model.Rows()
		columns := self.Model.Columns()

		for row := 0; row < rows; row++ {
			ctx.Response.XML.OpenTag("tr")
			if row&1 == 0 {
				ctx.Response.XML.Attrib("class", "row", row, " even")
			} else {
				ctx.Response.XML.Attrib("class", "row", row, " odd")
			}

			for col := 0; col < columns; col++ {
				if self.HeaderRow && row == 0 {
					ctx.Response.XML.OpenTag("th")
				} else {
					ctx.Response.XML.OpenTag("td")
				}
				if col&1 == 0 {
					ctx.Response.XML.Attrib("class", "col", col, " even")
				} else {
					ctx.Response.XML.Attrib("class", "col", col, " odd")
				}
				view, err := self.Model.CellView(row, col, ctx)
				if view != nil && err == nil {
					view.Init(view)
					err = view.Render(ctx)
				}
				if err != nil {
					return err
				}
				ctx.Response.XML.CloseTagAlways() // td/th
			}

			ctx.Response.XML.CloseTagAlways() // tr
		}
	}

	ctx.Response.XML.CloseTagAlways() // table
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// TableModel

type TableModel interface {
	Rows() int
	Columns() int
	CellView(row int, column int, ctx *Context) (view View, err error)
}

///////////////////////////////////////////////////////////////////////////////
// ViewsTableModel

type ViewsTableModel []Views

func (self ViewsTableModel) Rows() int {
	return len(self)
}

func (self ViewsTableModel) Columns() int {
	columns := 0
	for row := range self {
		x := len(self[row])
		if x > columns {
			columns = x
		}
	}
	return columns
}

func (self ViewsTableModel) CellView(row int, column int, ctx *Context) (view View, err error) {
	if column >= len(self[row]) {
		return nil, nil
	}
	return self[row][column], nil
}

///////////////////////////////////////////////////////////////////////////////
// StringsTableModel

type HTMLStringsTableModel [][]string

func (self HTMLStringsTableModel) Rows() int {
	return len(self)
}

func (self HTMLStringsTableModel) Columns() int {
	columns := 0
	for row := range self {
		x := len(self[row])
		if x > columns {
			columns = x
		}
	}
	return columns
}

func (self HTMLStringsTableModel) CellView(row int, column int, ctx *Context) (view View, err error) {
	return HTML(self[row][column]), nil
}

///////////////////////////////////////////////////////////////////////////////
// StringsTableModel

type EscapeStringsTableModel [][]string

func (self EscapeStringsTableModel) Rows() int {
	return len(self)
}

func (self EscapeStringsTableModel) Columns() int {
	columns := 0
	for row := range self {
		x := len(self[row])
		if x > columns {
			columns = x
		}
	}
	return columns
}

func (self EscapeStringsTableModel) CellView(row int, column int, ctx *Context) (view View, err error) {
	return Escape(self[row][column]), nil
}
