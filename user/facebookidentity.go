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

func (self *FacebookIdentity) URL(context *view.Context, args ...string) string {
	return self.ProfileURL()
}

func (self *FacebookIdentity) LinkContent(context *view.Context) view.View {
	return view.Escape(self.LinkTitle(context))
}

func (self *FacebookIdentity) LinkTitle(context *view.Context) string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *FacebookIdentity) LinkRel(context *view.Context) string {
	return ""
}
