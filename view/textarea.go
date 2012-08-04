package view

const (
	TextAreaDefaultCols = 80
	TextAreaDefaultRows = 3
)

///////////////////////////////////////////////////////////////////////////////
// TextArea

type TextArea struct {
	ViewBaseWithId
	Text        string
	Name        string
	Cols        int
	Rows        int
	Readonly    bool
	Disabled    bool
	TabIndex    int
	Class       string
	Placeholder string
}

func (self *TextArea) Render(response *Response) (err error) {
	response.XML.OpenTag("textarea").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)

	cols := self.Cols
	if cols == 0 {
		cols = TextAreaDefaultCols
	}
	rows := self.Rows
	if rows == 0 {
		rows = TextAreaDefaultRows
	}

	response.XML.Attrib("name", self.Name)
	response.XML.Attrib("rows", rows)
	response.XML.Attrib("cols", cols)
	response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.Readonly {
		response.XML.Attrib("readonly", "readonly")
	}
	if self.Disabled {
		response.XML.Attrib("disabled", "disabled")
	}
	response.XML.AttribIfNotDefault("placeholder", self.Placeholder)

	response.XML.EscapeContent(self.Text)

	response.XML.ForceCloseTag()
	return nil
}

//func (self *TextArea) SetText(text string) {
//	self.Text = text
//	ViewChanged(self)
//}
//
//func (self *TextArea) SetName(name string) {
//	self.Name = name
//	ViewChanged(self)
//}
//
//func (self *TextArea) SetCols(cols int) {
//	self.Cols = cols
//	ViewChanged(self)
//}
//
//func (self *TextArea) SetRows(rows int) {
//	self.Rows = rows
//	ViewChanged(self)
//}
//
//func (self *TextArea) SetReadonly(readonly bool) {
//	self.Readonly = readonly
//	ViewChanged(self)
//}
//
//func (self *TextArea) SetTabIndex(tabindex int) {
//	self.TabIndex = tabindex
//	ViewChanged(self)
//}
//
//func (self *TextArea) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
