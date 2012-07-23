package mongo

import (
	"github.com/ungerik/go-start/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterEqualCaseInsensitiveQuery

type filterEqualCaseInsensitiveQuery struct {
	filterQueryBase
	selector string
	str      string
}

func (self *filterEqualCaseInsensitiveQuery) bsonSelector() bson.M {
	s := escapeStringForRegex(self.str)
	return bson.M{self.selector: bson.RegEx{"^" + s + "$", "i"}}
}

func (self *filterEqualCaseInsensitiveQuery) Selector() string {
	return self.selector
}

//func (self *filterEqualCaseInsensitiveQuery) CompareString() string {
//	return self.str
//}
