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

func (self *SkypeIdentity) URL(args ...string) string {
	return self.CallURL()
}

func (self *SkypeIdentity) LinkContent(context *view.Context) view.View {
	return view.Escape(self.LinkTitle(context))
}

func (self *SkypeIdentity) LinkTitle(context *view.Context) string {
	return self.ID.Get()
}

func (self *SkypeIdentity) LinkRel(context *view.Context) string {
	return ""
}
