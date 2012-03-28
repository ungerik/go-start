package view

///////////////////////////////////////////////////////////////////////////////
// ViewBaseWithDynamicChildren

// ViewBaseWithDynamicChildren extends ViewBase with some helper methods
// to manage dynamic request-context specific child views.
type ViewBaseWithDynamicChildren struct {
	ViewBase
	children Views
}

func (self *ViewBaseWithDynamicChildren) Init(thisView View) {
	self.ViewBase.Init(thisView)
	self.children = Views{}
}

func (self *ViewBaseWithDynamicChildren) IterateChildren(callback IterateChildrenCallback) {
	for _, child := range self.children {
		if child != nil && !callback(self.thisView, child) {
			return
		}
	}
}

func (self *ViewBaseWithDynamicChildren) RemoveChildren() {
	for _, child := range self.children {
		if child != nil {
			child.OnRemove()
		}
	}
	self.children = self.children[0:0]
}

func (self *ViewBaseWithDynamicChildren) AddAndInitChild(child View) {
	self.children = append(self.children, child)
	child.Init(child)
}
