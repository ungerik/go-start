package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// Format

type Format struct {
	ViewBase
	Text   string
	Args   []interface{}
	Escape bool
}

func (self *Format) Render(context *Context, writer *utils.XMLWriter) (err error) {
	if self.Escape {
		writer.PrintfEscape(self.Text, self.Args...)
	} else {
		writer.Printf(self.Text, self.Args...)
	}
	return nil
}
