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
	if self == nil {
		return ""
	}
	if self.ID == "" {
		return ""
	}
	return "http://www.xing.com/profile/" + self.ID.Get()
}

func (self *XingIdentity) URL(ctx *view.Context) string {
	return self.ProfileURL()
}

func (self *XingIdentity) LinkContent(ctx *view.Context) view.View {
	if self == nil {
		return nil
	}
	return view.Escape(self.LinkTitle(ctx))
}

func (self *XingIdentity) LinkTitle(ctx *view.Context) string {
	if self == nil {
		return ""
	}
	return self.ID.Get()
}

func (self *XingIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
