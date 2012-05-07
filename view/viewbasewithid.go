package view

///////////////////////////////////////////////////////////////////////////////
// ViewBaseWithId

// ViewBaseWithId extends ViewBase with an id for the view.
type ViewBaseWithId struct {
	ViewBase
	id string
}

func (self *ViewBaseWithId) Init(thisView View) {
	self.ViewBase.Init(thisView)
	self.id = NewViewID(thisView)
}

func (self *ViewBaseWithId) ID() string {
	return self.id
}
