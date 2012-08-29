package user

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// GitHubIdentity

type GitHubIdentity struct {
	ID          model.String
	Name        model.String
	Confirmed   model.DateTime
	AccessToken model.String
}

func (self *GitHubIdentity) ProfileURL() string {
	if self.Name == "" {
		return ""
	}
	return "https://github.com/" + self.Name.Get()
}

func (self *GitHubIdentity) URL(ctx *view.Context) string {
	return self.ProfileURL()
}

func (self *GitHubIdentity) LinkContent(ctx *view.Context) view.View {
	return view.Escape(self.LinkTitle(ctx))
}

func (self *GitHubIdentity) LinkTitle(ctx *view.Context) string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *GitHubIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
