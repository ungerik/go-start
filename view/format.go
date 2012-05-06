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

func (self *Format) Render(response *Response) (err error) {
	writer := utils.NewXMLWriter(response)
	if self.Escape {
		writer.PrintfEscape(self.Text, self.Args...)
	} else {
		writer.Printf(self.Text, self.Args...)
	}
	return nil
}
