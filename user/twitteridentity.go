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

func (self *TwitterIdentity) URL(args ...string) string {
	return self.ProfileURL()
}

func (self *TwitterIdentity) LinkContent(urlArgs ...string) view.View {
	return view.Escape(self.LinkTitle(urlArgs...))
}

func (self *TwitterIdentity) LinkTitle(urlArgs ...string) string {
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

func (self *TwitterIdentity) LinkRel() string {
	return ""
}
