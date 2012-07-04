package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterContainsQuery

type filterContainsQuery struct {
	filterQueryBase
	selector        string
	str             string
	caseInsensitive bool
}

func (self *filterContainsQuery) bsonSelector() bson.M {
	s := escapeStringForRegex(self.str)
	var options string
	if self.caseInsensitive {
		options = "i"
	}
	return bson.M{self.selector: bson.RegEx{s, options}}
}

func (self *filterContainsQuery) Selector() string {
	return self.selector
}

//func (self *filterContainsQuery) CompareString() string {
//	return self.str
//}
//
//func (self *filterContainsQuery) CaseInsensitive() bool {
//	return self.caseInsensitive
//}

