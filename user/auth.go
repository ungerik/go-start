package user

import "github.com/ungerik/go-start/view"

///////////////////////////////////////////////////////////////////////////////
// Auth

type Auth struct {
	LoginURL view.URL
}

func (self *Auth) Authenticate(context *view.Context) (ok bool, err error) {
	id, ok := context.SessionID()
	if !ok {
		return false, nil
	}

	ok, err = IsConfirmedUserID(id)
	if !ok && err == nil && self.LoginURL != nil {
		err = view.Redirect(self.LoginURL.URL(context.PathArgs...))
	}
	return ok, err
}
