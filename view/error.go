package view

func NewError(err error) *Error {
	return &Error{Err: err}
}

// Error wraps an error and returns it at the Render() method.
type Error struct {
	Err error
}

func (self *Error) Init(thisView View) {
}

func (self *Error) ID() string {
	return ""
}

func (self *Error) IterateChildren(callback IterateChildrenCallback) {
}

func (self *Error) Render(ctx *Context) (err error) {
	return self.Err
}

func (self *Error) Error() string {
	return self.Err.Error()
}
