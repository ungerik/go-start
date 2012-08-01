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

func (self *SkypeIdentity) LinkContent(urlArgs ...string) view.View {
	return view.Escape(self.LinkTitle(urlArgs...))
}

func (self *SkypeIdentity) LinkTitle(urlArgs ...string) string {
	return self.ID.Get()
}

func (self *SkypeIdentity) LinkRel() string {
	return ""
}
