package user

import (
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// LoginNav

type LoginNav struct {
	view.ViewBaseWithId
	Class string

	// The Links won't be enumerated and initialized as child views
	LoginSignup *view.Link
	Logout      *view.Link
}

func (self *LoginNav) Render(context *view.Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("div").Attrib("id", self.ID()).AttribIfNotDefault("class", self.Class)
	if _, ok := context.SessionID(); ok {
		err = self.Logout.Render(context, writer)
	} else {
		err = self.LoginSignup.Render(context, writer)
	}
	writer.CloseTag()
	return err
}
