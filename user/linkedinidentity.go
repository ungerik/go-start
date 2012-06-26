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
	if !self.Name.IsEmpty() {
		name := self.Name.Get()
		if strings.Index(name, "/") == -1 {
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

func (self *LinkedInIdentity) URL(args ...string) string {
	return self.ProfileURL()
}

func (self *LinkedInIdentity) LinkContent(context *view.Context) view.View {
	return view.Escape(self.LinkTitle(context))
}

func (self *LinkedInIdentity) LinkTitle(context *view.Context) string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *LinkedInIdentity) LinkRel(context *view.Context) string {
	return ""
}
