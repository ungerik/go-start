package view

///////////////////////////////////////////////////////////////////////////////
// Format

type Format struct {
	ViewBase
	Text   string
	Args   []interface{}
	Escape bool
}

func (self *Format) Render(response *Response) (err error) {
	if self.Escape {
		response.XML.PrintfEscape(self.Text, self.Args...)
	} else {
		response.XML.Printf(self.Text, self.Args...)
	}
	return nil
}
