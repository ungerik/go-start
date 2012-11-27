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
	Required    bool // HTML5
	Autofocus   bool // HTML5
}

func (self *TextArea) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("textarea")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)

	cols := self.Cols
	if cols == 0 {
		cols = TextAreaDefaultCols
	}
	rows := self.Rows
	if rows == 0 {
		rows = TextAreaDefaultRows
	}

	ctx.Response.XML.Attrib("name", self.Name)
	ctx.Response.XML.Attrib("rows", rows)
	ctx.Response.XML.Attrib("cols", cols)
	ctx.Response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	ctx.Response.XML.AttribFlag("readonly", self.Readonly)
	ctx.Response.XML.AttribFlag("disabled", self.Disabled)
	ctx.Response.XML.AttribFlag("required", self.Required)
	ctx.Response.XML.AttribFlag("autofocus", self.Autofocus)
	ctx.Response.XML.AttribIfNotDefault("placeholder", self.Placeholder)

	ctx.Response.XML.EscapeContent(self.Text)

	ctx.Response.XML.CloseTagAlways()
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
