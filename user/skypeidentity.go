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

func (self *SkypeIdentity) URL(ctx *view.Context) string {
	return self.CallURL()
}

func (self *SkypeIdentity) LinkContent(ctx *view.Context) view.View {
	return view.Escape(self.LinkTitle(ctx))
}

func (self *SkypeIdentity) LinkTitle(ctx *view.Context) string {
	return self.ID.Get()
}

func (self *SkypeIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
