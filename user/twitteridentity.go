package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// TwitterIdentity

type TwitterIdentity struct {
	ID          model.String
	Name        model.String
	Confirmed   model.DateTime
	AccessToken model.String
}

func (self *TwitterIdentity) ProfileURL() string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	if name[0] == '@' {
		name = name[1:]
		if name == "" {
			return ""
		}
	}
	return "http://twitter.com/" + name
}

func (self *TwitterIdentity) URL(request *view.Request, session *view.Session, response *view.Response, args ...string) string {
	return self.ProfileURL()
}

func (self *TwitterIdentity) LinkContent(request *view.Request, session *view.Session, response *view.Response) view.View {
	return view.Escape(self.LinkTitle(context))
}

func (self *TwitterIdentity) LinkTitle(request *view.Request, session *view.Session, response *view.Response) string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	if name[0] != '@' {
		name = "@" + name
	}
	return name
}

func (self *TwitterIdentity) LinkRel(request *view.Request, session *view.Session, response *view.Response) string {
	return ""
}
