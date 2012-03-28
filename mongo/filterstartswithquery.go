package mongo

import (
	"launchpad.net/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterStartsWithQuery

type filterStartsWithQuery struct {
	filterQueryBase
	selector        string
	str             string
	caseInsensitive bool
}

func (self *filterStartsWithQuery) bsonSelector() bson.M {
	s := escapeStringForRegex(self.str)
	var options string
	if self.caseInsensitive {
		options = "i"
	}
	return bson.M{self.selector: bson.RegEx{"^" + s, options}}
}

func (self *filterStartsWithQuery) Selector() string {
	return self.selector
}
