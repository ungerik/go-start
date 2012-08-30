package view

///////////////////////////////////////////////////////////////////////////////
// Format

type Format struct {
	ViewBase
	Text   string
	Args   []interface{}
	Escape bool
}

func (self *Format) Render(ctx *Context) (err error) {
	if self.Escape {
		ctx.Response.XML.PrintfEscape(self.Text, self.Args...)
	} else {
		ctx.Response.XML.Printf(self.Text, self.Args...)
	}
	return nil
}
