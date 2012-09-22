package view

type Comment string

func (self Comment) Init(thisView View) {
}

func (self Comment) ID() string {
	return ""
}

func (self Comment) IterateChildren(callback IterateChildrenCallback) {
}

func (self Comment) Render(ctx *Context) (err error) {
	ctx.Response.Printf("<!-- %s -->", self)
	return nil
}
