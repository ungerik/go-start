package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// XingIdentity

type XingIdentity struct {
	ID model.String
}

func (self *XingIdentity) ProfileURL() string {
	if self.ID == "" {
		return ""
	}
	return "http://www.xing.com/profile/" + self.ID.Get()
}

func (self *XingIdentity) URL(request *view.Request, session *view.Session, response *view.Response, args ...string) string {
	return self.ProfileURL()
}

func (self *XingIdentity) LinkContent(request *view.Request, session *view.Session, response *view.Response) view.View {
	return view.Escape(self.LinkTitle(context))
}

func (self *XingIdentity) LinkTitle(request *view.Request, session *view.Session, response *view.Response) string {
	return self.ID.Get()
}

func (self *XingIdentity) LinkRel(request *view.Request, session *view.Session, response *view.Response) string {
	return ""
}
