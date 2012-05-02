package view

import "github.com/ungerik/go-start/utils"

const (
	TextAreaDefaultCols = 80
	TextAreaDefaultRows = 3
)

///////////////////////////////////////////////////////////////////////////////
// TextArea

type TextArea struct {
	ViewBaseWithId
	Text     string
	Name     string
	Cols     int
	Rows     int
	Readonly bool
	Disabled bool
	TabIndex int
	Class    string
}

func (self *TextArea) Render(request *Request, session *Session, response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	writer.OpenTag("textarea").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)

	cols := self.Cols
	if cols == 0 {
		cols = TextAreaDefaultCols
	}
	rows := self.Rows
	if rows == 0 {
		rows = TextAreaDefaultRows
	}

	writer.Attrib("name", self.Name)
	writer.Attrib("rows", rows)
	writer.Attrib("cols", cols)
	writer.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.Readonly {
		writer.Attrib("readonly", "readonly")
	}
	if self.Disabled {
		writer.Attrib("disabled", "disabled")
	}

	writer.EscapeContent(self.Text)
	writer.CloseTag()
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
