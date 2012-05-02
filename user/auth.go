package user

import "github.com/ungerik/go-start/view"

///////////////////////////////////////////////////////////////////////////////
// Auth

type Auth struct {
	LoginURL view.URL
}

func (self *Auth) Authenticate(request *view.Request, session *view.Session, response *view.Response) (ok bool, err error) {
	id, ok := session.ID()
	if !ok {
		return false, nil
	}

	ok, err = IsConfirmedUserID(id)
	if !ok && err == nil && self.LoginURL != nil {
		err = view.Redirect(self.LoginURL.URL(request, session, response))
	}
	return ok, err
}
