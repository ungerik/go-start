package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// XingIdentity

type XingIdentity struct {
	ID model.String
}

func (self *XingIdentity) ProfileURL() string {
	if self.ID == "" {
		return ""
	}
	return "http://www.xing.com/profile/" + self.ID.Get()
}

func (self *XingIdentity) URL(response *view.Response, args ...string) string {
	return self.ProfileURL()
}

func (self *XingIdentity) LinkContent(response *view.Response) view.View {
	return view.Escape(self.LinkTitle(response))
}

func (self *XingIdentity) LinkTitle(response *view.Response) string {
	return self.ID.Get()
}

func (self *XingIdentity) LinkRel(response *view.Response) string {
	return ""
}
