package view

import (
	"github.com/ungerik/go-start/debug"
)

const (
	RichTextAreaDefaultCols = 80
	RichTextAreaDefaultRows = 10
)

///////////////////////////////////////////////////////////////////////////////
// TextArea

type RichTextArea struct {
	ViewBaseWithId
	Text                 string
	Name                 string
	Cols                 int
	Rows                 int
	Readonly             bool
	Disabled             bool
	TabIndex             int
	Class                string
	Placeholder          string
	ToolbarStylesheetURL string
}

// We use https://github.com/xing/wysihtml5/ as open source rich text editor
//
func (self *RichTextArea) Render(ctx *Context) (err error) {

	if !Config.RichText.UseGlobalCSS {
		if self.ToolbarStylesheetURL != "" {
			debug.Print("Loading StylesheetURL")
			ctx.Response.RequireStyleURL(self.ToolbarStylesheetURL, 0)
		} else {
			debug.Print("Loading Config")
			ctx.Response.RequireStyleURL(Config.RichText.ToolbarCSS, 0)
		}
	}
	// ctx.Response.RequireStyleURL("http://yui.yahooapis.com/2.9.0/build/reset/reset-min.css", 0)

	// wysihtml5 parser rules
	ctx.Response.RequireScriptURL("/js/libs/wysihtml5-advanced.js", 0)
	// wysihtml5 Library
	ctx.Response.RequireScriptURL("/js/libs/wysihtml5-0.3.0.min.js", 1)

	ctx.Response.RequireScript(`var editor = new wysihtml5.Editor("`+self.id+`", { // id of textarea element
  toolbar:      "wysihtml5-editor-toolbar", // id of toolbar element
  parserRules:  wysihtml5ParserRules, // defined in parser rules set 
  composerClassName:    "wysihtml5-editor",
  stylesheets: ["`+Config.RichText.EditorCSS+`"]
});`, 2)

	// ctx.Response.XML.Content(`<div id="wysihtml5-toolbar" style="display: none;">` + self.ToolbarHtml + `</div>`)

	ctx.Response.XML.Content(`<div id="wysihtml5-editor-toolbar" style="display: none;">` + Config.RichText.DefaultToolbar + `</div>`)
	// ctx.Response.XML.OpenTag("form")
	ctx.Response.XML.OpenTag("textarea")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)

	ctx.Response.XML.Attrib("class", "wysihtml5 "+self.Class)

	cols := self.Cols
	if cols == 0 {
		cols = RichTextAreaDefaultCols
	}
	rows := self.Rows
	if rows == 0 {
		rows = RichTextAreaDefaultRows
	}

	ctx.Response.XML.Attrib("name", self.Name)
	ctx.Response.XML.Attrib("rows", rows)
	ctx.Response.XML.Attrib("cols", cols)
	ctx.Response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	ctx.Response.XML.AttribFlag("readonly", self.Readonly)
	ctx.Response.XML.AttribFlag("disabled", self.Disabled)
	ctx.Response.XML.AttribIfNotDefault("placeholder", self.Placeholder)

	ctx.Response.XML.Attrib("autofocus", "autofocus")

	ctx.Response.XML.EscapeContent(self.Text)

	ctx.Response.XML.CloseTagAlways()
	// ctx.Response.XML.CloseTagAlways()
	return nil
}

func SetDefaultToolbar() string {
	return `
  <a data-wysihtml5-command="bold">bold</a>
  <a data-wysihtml5-command="italic">italic</a>
  
  <!-- Some wysihtml5 commands require extra parameters -->
  <a data-wysihtml5-command="foreColor" data-wysihtml5-command-value="red">red</a>
  <a data-wysihtml5-command="foreColor" data-wysihtml5-command-value="green">green</a>
  <a data-wysihtml5-command="foreColor" data-wysihtml5-command-value="blue">blue</a>
  
  <!-- Some wysihtml5 commands like 'createLink' require extra paramaters specified by the user (eg. href) -->
  <a data-wysihtml5-command="createLink">insert link</a>
  <div data-wysihtml5-dialog="createLink" style="display: none;">
    <label>
      Link:
      <input data-wysihtml5-dialog-field="href" value="http://" class="text">
    </label>
    <a data-wysihtml5-dialog-action="save">OK</a> <a data-wysihtml5-dialog-action="cancel">Cancel</a>
  </div>`
}

//func (self *RichTextArea) SetText(text string) {
//	self.Text = text
//	ViewChanged(self)
//}
//
//func (self *RichTextArea) SetName(name string) {
//	self.Name = name
//	ViewChanged(self)
//}
//
//func (self *RichTextArea) SetCols(cols int) {
//	self.Cols = cols
//	ViewChanged(self)
//}
//
//func (self *RichTextArea) SetRows(rows int) {
//	self.Rows = rows
//	ViewChanged(self)
//}
//
//func (self *RichTextArea) SetReadonly(readonly bool) {
//	self.Readonly = readonly
//	ViewChanged(self)
//}
//
//func (self *RichTextArea) SetTabIndex(tabindex int) {
//	self.TabIndex = tabindex
//	ViewChanged(self)
//}
//
//func (self *RichTextArea) SetClass(class string) {
//	self.Class = class
//	ViewChanged(self)
//}
