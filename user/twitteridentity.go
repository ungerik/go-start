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

func (self *TwitterIdentity) URL(ctx *view.Context) string {
	return self.ProfileURL()
}

func (self *TwitterIdentity) LinkContent(ctx *view.Context) view.View {
	return view.Escape(self.LinkTitle(ctx))
}

func (self *TwitterIdentity) LinkTitle(ctx *view.Context) string {
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

func (self *TwitterIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
