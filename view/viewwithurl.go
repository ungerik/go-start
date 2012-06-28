package view

// ViewWithURL combines the View interface with the URL interface
// for views that have an URL
type ViewWithURL interface {
	View
	URL
	SetPath(path string)
}
