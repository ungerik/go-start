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

func (self *SkypeIdentity) URL(response *view.Response, args ...string) string {
	return self.CallURL()
}

func (self *SkypeIdentity) LinkContent(response *view.Response) view.View {
	return view.Escape(self.LinkTitle(response))
}

func (self *SkypeIdentity) LinkTitle(response *view.Response) string {
	return self.ID.Get()
}

func (self *SkypeIdentity) LinkRel(response *view.Response) string {
	return ""
}
