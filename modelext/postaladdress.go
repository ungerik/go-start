package modelext

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// PostalAddress

type PostalAddress struct {
	FirstLine  model.String
	SecondLine model.String
	ZIP        model.String
	City       model.String
	State      model.String
	Country    model.Country
}

func (self *PostalAddress) String() string {
	return self.StringSep(", ")
}

func (self *PostalAddress) StringSep(sep string) string {
	return utils.JoinNonEmptyStrings(
		sep, self.FirstLine.Get(),
		self.SecondLine.Get(),
		self.ZIP.Get(),
		self.City.Get(),
		self.State.Get(),
		self.Country.Get(),
	)
}
