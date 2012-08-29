package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// FacebookIdentity

type FacebookIdentity struct {
	ID          model.String
	Name        model.String `view:"placeholder=your facebook url"`
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

func (self *FacebookIdentity) URL(ctx *view.Context) string {
	return self.ProfileURL()
}

func (self *FacebookIdentity) LinkContent(ctx *view.Context) view.View {
	return view.Escape(self.LinkTitle(ctx))
}

func (self *FacebookIdentity) LinkTitle(ctx *view.Context) string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *FacebookIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
