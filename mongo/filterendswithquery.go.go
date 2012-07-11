package mongo

import (
	"github.com/ungerik/go-start/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterEndsWithQuery

type filterEndsWithQuery struct {
	filterQueryBase
	selector        string
	str             string
	caseInsensitive bool
}

func (self *filterEndsWithQuery) bsonSelector() bson.M {
	s := escapeStringForRegex(self.str)
	var options string
	if self.caseInsensitive {
		options = "i"
	}
	return bson.M{self.selector: bson.RegEx{s + "$", options}}
}

func (self *filterEndsWithQuery) Selector() string {
	return self.selector
}

//func (self *filterEndsWithQuery) CompareString() string {
//	return self.str
//}
//
//func (self *filterEndsWithQuery) CaseInsensitive() bool {
//	return self.caseInsensitive
//}
