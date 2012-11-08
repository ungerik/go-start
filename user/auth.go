package user

import "github.com/ungerik/go-start/view"

///////////////////////////////////////////////////////////////////////////////
// Auth

type Auth struct {
	LoginURL view.URL
}

func (self *Auth) Authenticate(ctx *view.Context) (ok bool, err error) {
	id := ctx.Session.ID()
	if id == "" {
		return false, nil
	}

	/// todo: Use hashed logged in cookie instead of user ID!!
	ok, err = IsConfirmedUserID(id)
	if !ok && err == nil && self.LoginURL != nil {
		err = view.Redirect(self.LoginURL.URL(ctx))
	}
	return ok, err
}
