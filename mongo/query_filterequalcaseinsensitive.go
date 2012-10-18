package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterEqualCaseInsensitive

type query_filterEqualCaseInsensitive struct {
	query_filterBase
	selector string
	str      string
}

func (self *query_filterEqualCaseInsensitive) bsonSelector() bson.M {
	s := escapeStringForRegex(self.str)
	return bson.M{self.selector: bson.RegEx{"^" + s + "$", "i"}}
}

func (self *query_filterEqualCaseInsensitive) Selector() string {
	return self.selector
}

//func (self *query_filterEqualCaseInsensitive) CompareString() string {
//	return self.str
//}
