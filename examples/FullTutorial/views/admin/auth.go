package admin

import (
	"github.com/ungerik/go-start/user"
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

type Admin_Authenticator struct{}

func (self *Admin_Authenticator) Authenticate(ctx *Context) (ok bool, err error) {
	var u models.User
	found, err := user.OfSession(ctx.Session, &u)
	if !found {
		return false, err
	}
	return u.Admin.Get(), nil
}

func init() {
	Admin_Auth = new(Admin_Authenticator)
}
