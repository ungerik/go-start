package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// SkypeIdentity

type SkypeIdentity struct {
	ID        model.String
	Confirmed model.DateTime
}

func (self *SkypeIdentity) CallURL() string {
	if self.ID == "" {
		return ""
	}
	return "skype:" + self.ID.Get()
}

func (self *SkypeIdentity) URL(request *view.Request, session *view.Session, response *view.Response, args ...string) string {
	return self.CallURL()
}

func (self *SkypeIdentity) LinkContent(request *view.Request, session *view.Session, response *view.Response) view.View {
	return view.Escape(self.LinkTitle(request, session, response))
}

func (self *SkypeIdentity) LinkTitle(request *view.Request, session *view.Session, response *view.Response) string {
	return self.ID.Get()
}

func (self *SkypeIdentity) LinkRel(request *view.Request, session *view.Session, response *view.Response) string {
	return ""
}
