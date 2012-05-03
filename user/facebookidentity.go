package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// FacebookIdentity

type FacebookIdentity struct {
	ID          model.String
	Name        model.String
	Confirmed   model.DateTime
	AccessToken model.String
}

func (self *FacebookIdentity) ProfileURL() string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return "http://facebook.com/" + name
}

func (self *FacebookIdentity) URL(request *view.Request, session *view.Session, response *view.Response, args ...string) string {
	return self.ProfileURL()
}

func (self *FacebookIdentity) LinkContent(request *view.Request, session *view.Session, response *view.Response) view.View {
	return view.Escape(self.LinkTitle(request, session, response))
}

func (self *FacebookIdentity) LinkTitle(request *view.Request, session *view.Session, response *view.Response) string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *FacebookIdentity) LinkRel(request *view.Request, session *view.Session, response *view.Response) string {
	return ""
}
