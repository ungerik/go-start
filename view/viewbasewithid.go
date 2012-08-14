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
}

func (self *ViewBaseWithId) ID() string {
	if self.id == "" {
		self.id = NewViewID(self.thisView)
	}
	return self.id
}
