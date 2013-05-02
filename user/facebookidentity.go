package user

import (
	"fmt"

	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

///////////////////////////////////////////////////////////////////////////////
// FacebookIdentity

type FacebookIdentity struct {
	ID          model.String
	Name        model.String `view:"placeholder=your facebook url"`
	Confirmed   model.DateTime
	AccessToken model.String
}

func (self *FacebookIdentity) NameOrID() string {
	if self == nil {
		return ""
	}
	if !self.Name.IsEmpty() {
		return self.Name.Get()
	}
	return self.ID.Get()
}

func (self *FacebookIdentity) ProfileURL() string {
	if self == nil {
		return ""
	}
	return "http://facebook.com/" + self.NameOrID()
}

func (self *FacebookIdentity) ProfileImageURL() string {
	name := self.NameOrID()
	if name == "" {
		return ""
	}
	return fmt.Sprintf("http://graph.facebook.com/%s/picture?type=large", name)
}

func (self *FacebookIdentity) URL(ctx *view.Context) string {
	return self.ProfileURL()
}

func (self *FacebookIdentity) LinkContent(ctx *view.Context) view.View {
	if self == nil {
		return nil
	}
	return view.Escape(self.LinkTitle(ctx))
}

func (self *FacebookIdentity) LinkTitle(ctx *view.Context) string {
	return self.NameOrID()
}

func (self *FacebookIdentity) LinkRel(ctx *view.Context) string {
	return ""
}
