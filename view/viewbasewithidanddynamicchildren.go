package view

///////////////////////////////////////////////////////////////////////////////
// ViewBaseWithIdAndDynamicChildren

// ViewBaseWithDynamicChildren extends ViewBaseWithId with some helper methods
// to manage dynamic request-context specific child views.
type ViewBaseWithIdAndDynamicChildren struct {
	ViewBaseWithId
	children Views
}

func (self *ViewBaseWithIdAndDynamicChildren) Init(thisView View) {
	self.ViewBaseWithId.Init(thisView)
	self.children = Views{}
}

func (self *ViewBaseWithIdAndDynamicChildren) IterateChildren(callback IterateChildrenCallback) {
	for _, child := range self.children {
		if child != nil && !callback(self.thisView, child) {
			return
		}
	}
}

func (self *ViewBaseWithIdAndDynamicChildren) RemoveChildren() {
	for _, child := range self.children {
		if child != nil {
			child.OnRemove()
		}
	}
	self.children = self.children[0:0]
}

func (self *ViewBaseWithIdAndDynamicChildren) AddAndInitChild(child View) {
	self.children = append(self.children, child)
	child.Init(child)
}
