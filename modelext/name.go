package modelext

import (
	"strings"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// Name

type Name struct {
	Prefix       model.String
	First        model.String
	Middle       model.String
	Last         model.String
	Postfix      model.String
	Organization model.String
}

func (self *Name) SetForPerson(prefix, first, middle, last, postfix string) {
	self.Prefix.Set(prefix)
	self.First.Set(strings.Title(strings.ToLower(first)))
	self.Middle.Set(strings.Title(strings.ToLower(middle)))
	self.Last.Set(strings.Title(strings.ToLower(last)))
	self.Postfix.Set(postfix)
	self.Organization = ""
}

func (self *Name) SetForOrganization(organization string) {
	self.Prefix = ""
	self.First = ""
	self.Middle = ""
	self.Last = ""
	self.Postfix = ""
	self.Organization.Set(organization)
}

func (self *Name) String() string {
	if self.Organization != "" {
		return self.Organization.Get()
	} 
	return utils.JoinNonEmptyStrings(" ", self.Prefix.Get(), self.First.Get(), self.Middle.Get(), self.Last.Get(), self.Postfix.Get())
}
