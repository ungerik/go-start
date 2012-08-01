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

func (self *FacebookIdentity) URL(args ...string) string {
	return self.ProfileURL()
}

func (self *FacebookIdentity) LinkContent(urlArgs ...string) view.View {
	return view.Escape(self.LinkTitle(urlArgs...))
}

func (self *FacebookIdentity) LinkTitle(urlArgs ...string) string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *FacebookIdentity) LinkRel() string {
	return ""
}
