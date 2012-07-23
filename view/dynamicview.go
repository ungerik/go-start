package view

/*
DynamicView implements View for a function that creates and renders a dynamic
child-view in the Render method.

Example:

	dynamicView := DynamicView(
		func(response *Response) (view View, err error) {
			return HTML("return dynamic created views here"), nil
		},
	)
*/
type DynamicView func(response *Response) (view View, err error)

func (self DynamicView) Init(thisView View) {
}

func (self DynamicView) ID() string {
	return ""
}

func (self DynamicView) IterateChildren(callback IterateChildrenCallback) {
}

func (self DynamicView) Render(response *Response) error {
	child, err := self(response)
	if err != nil || child == nil {
		return err
	}
	child.Init(child)
	return child.Render(response)
}
