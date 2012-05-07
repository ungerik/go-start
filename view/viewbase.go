package view

// import "github.com/ungerik/go-start/debug"

///////////////////////////////////////////////////////////////////////////////
// ViewBase

// ViewBase is a base for View implementations.
type ViewBase struct {
	// thisView holds the implemented View as seen from outside of ViewBase
	// because ViewBase can't know how View will be implemented.
	thisView View
}

func (self *ViewBase) Init(thisView View) {
	if thisView == self.thisView {
		return // already initialized
	}
	self.thisView = thisView
	thisView.IterateChildren(func(parent View, child View) (next bool) {
		if child != nil {
			child.Init(child)
		}
		return true
	})
}

func (self *ViewBase) ID() string {
	return ""
}

func (self *ViewBase) IterateChildren(callback IterateChildrenCallback) {
}
