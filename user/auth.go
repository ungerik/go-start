package user

import "github.com/ungerik/go-start/view"

///////////////////////////////////////////////////////////////////////////////
// Auth

type Auth struct {
	LoginURL view.URL
}

func (self *Auth) Authenticate(response *view.Response) (ok bool, err error) {
	id, ok := response.Session.ID()
	if !ok {
		return false, nil
	}

	ok, err = IsConfirmedUserID(id)
	if !ok && err == nil && self.LoginURL != nil {
		err = view.Redirect(self.LoginURL.URL(response.Request.URLArgs...))
	}
	return ok, err
}
