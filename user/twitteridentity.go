package user

import (
	"fmt"

	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// TwitterIdentity

type TwitterIdentity struct {
	ID          model.String
	Name        model.String
	Confirmed   model.DateTime
	AccessToken model.String
}

func (self *TwitterIdentity) NameOrID() string {
	name := self.Name.Get()
	if name == "" {
		name = self.ID.Get()
		if name == "" {
			return ""
		}
	}
	if name[0] == '@' {
		name = name[1:]
		if name == "" {
			return ""
		}
	}
	return name
}

func (self *TwitterIdentity) ProfileURL() string {
	if self == nil {
		return ""
	}
	return "http://twitter.com/" + self.NameOrID()
}

func (self *TwitterIdentity) ProfileImageURL() string {
	if self == nil {
		return ""
	}
	name := self.NameOrID()
	if name == "" {
		return ""
	}
	return fmt.Sprintf("https://api.twitter.com/1/users/profile_image/%s?size=bigger", name)
}

func (self *TwitterIdentity) URL(ctx *view.Context) string {
	return self.ProfileURL()
}

func (self *TwitterIdentity) LinkContent(ctx *view.Context) view.View {
	if self == nil {
		return nil
	}
	return view.Escape(self.LinkTitle(ctx))
}

func (self *TwitterIdentity) LinkTitle(ctx *view.Context) string {
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
	if name[0] != '@' {
		name = "@" + name
	}
	return name
}

func (self *TwitterIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
