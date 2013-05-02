package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// LinkedInIdentity

type LinkedInIdentity struct {
	ID          model.String
	Name        model.String
	Confirmed   model.DateTime
	AccessToken model.String
}

func (self *LinkedInIdentity) ProfileURL() string {
	if self == nil {
		return ""
	}
	if !self.Name.IsEmpty() {
		name := self.Name.Get()
		if strings.IndexRune(name, '/') == -1 {
			return "http://linkedin.com/in/" + name
		} else {
			return "http://linkedin.com/pub/" + name
		}
	}
	if !self.ID.IsEmpty() {
		return "http://linkedin.com/profile/view?id=" + self.ID.Get()
	}
	return ""
}

func (self *LinkedInIdentity) URL(ctx *view.Context) string {
	return self.ProfileURL()
}

func (self *LinkedInIdentity) LinkContent(ctx *view.Context) view.View {
	if self == nil {
		return nil
	}
	return view.Escape(self.LinkTitle(ctx))
}

func (self *LinkedInIdentity) LinkTitle(ctx *view.Context) string {
	if self == nil {
		return ""
	}
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *LinkedInIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
