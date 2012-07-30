package modelext

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// PostalAddress

type PostalAddress struct {
	FirstLine  model.String `view:"size=40"`
	SecondLine model.String `view:"size=40"`
	ZIP        model.String `view:"size=10"`
	City       model.String `view:"size=20"`
	State      model.String `view:"size=20"`
	Country    model.Country
}

func (self *PostalAddress) String() string {
	return self.StringSep(", ")
}

func (self *PostalAddress) StringSep(sep string) string {
	return utils.JoinNonEmptyStrings(
		sep, self.FirstLine.Get(),
		self.SecondLine.Get(),
		utils.JoinNonEmptyStrings(" ", self.ZIP.Get(), self.City.Get()),
		self.State.Get(),
		self.Country.Get(),
	)
}
