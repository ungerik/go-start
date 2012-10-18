package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterContains

type query_filterContains struct {
	query_filterBase
	selector        string
	str             string
	caseInsensitive bool
}

func (self *query_filterContains) bsonSelector() bson.M {
	s := escapeStringForRegex(self.str)
	var options string
	if self.caseInsensitive {
		options = "i"
	}
	return bson.M{self.selector: bson.RegEx{s, options}}
}

func (self *query_filterContains) Selector() string {
	return self.selector
}

//func (self *query_filterContains) CompareString() string {
//	return self.str
//}
//
//func (self *query_filterContains) CaseInsensitive() bool {
//	return self.caseInsensitive
//}
