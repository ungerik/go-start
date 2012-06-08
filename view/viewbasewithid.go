package view

///////////////////////////////////////////////////////////////////////////////
// ViewBaseWithId

// ViewBaseWithId extends ViewBase with an id for the view.
type ViewBaseWithId struct {
	ViewBase
	id string
}

func (self *ViewBaseWithId) Init(thisView View) {
	if thisView == self.thisView {
		return // already initialized
	}
	self.ViewBase.Init(thisView)
	self.id = NewViewID(thisView)
}

func (self *ViewBaseWithId) ID() string {
	return self.id
}
