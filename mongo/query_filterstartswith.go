package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterStartsWith

type query_filterStartsWith struct {
	query_filterBase
	selector        string
	str             string
	caseInsensitive bool
}

func (self *query_filterStartsWith) bsonSelector() bson.M {
	s := escapeStringForRegex(self.str)
	var options string
	if self.caseInsensitive {
		options = "i"
	}
	return bson.M{self.selector: bson.RegEx{"^" + s, options}}
}

func (self *query_filterStartsWith) Selector() string {
	return self.selector
}
