package view

/*
DynamicView implements View for a function that creates and renders a dynamic
child-view in the Render method.

Example:

	dynamicView := DynamicView(
		func(request *Request, session *Session, response *Response) (view View, err error) {
			return HTML("return dynamic created views here"), nil
		},
	)
*/
type DynamicView func(request *Request, session *Session, response *Response) (view View, err error)

func (self DynamicView) Init(thisView View) {
}

func (self DynamicView) OnRemove() {
}

func (self DynamicView) ID() string {
	return ""
}

func (self DynamicView) IterateChildren(callback IterateChildrenCallback) {
}

func (self DynamicView) Render(request *Request, session *Session, response *Response) error {
	child, err := self(request, session, response)
	if err != nil || child == nil {
		return err
	}
	child.Init(child)
	return child.Render(request, session, response)
}
