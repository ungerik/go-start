package view

import "github.com/ungerik/go-start/utils"

///////////////////////////////////////////////////////////////////////////////
// RedirectView

// If rendered, this view will cause a HTTP redirect.
type RedirectView struct {
	ViewBase
	URL       string
	Permanent bool
}

func (self *RedirectView) Render(context *Context, writer *utils.XMLWriter) (err error) {
	if self.Permanent {
		return PermanentRedirect(self.URL)
	}
	return Redirect(self.URL)
}
