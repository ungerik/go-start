package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterLess

type query_filterLess struct {
	query_filterBase
	selector string
	value    interface{}
}

func (self *query_filterLess) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$lt": self.value}}
}

func (self *query_filterLess) Selector() string {
	return self.selector
}
