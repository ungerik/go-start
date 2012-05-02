package view

///////////////////////////////////////////////////////////////////////////////
// HTML

type HTML string

func (self HTML) Init(thisView View) {
}

func (self HTML) OnRemove() {
}

func (self HTML) ID() string {
	return ""
}

func (self HTML) IterateChildren(callback IterateChildrenCallback) {
}

func (self HTML) Render(request *Request, session *Session, response *Response) (err error) {
	response.WriteString(string(self))
	return nil
}
