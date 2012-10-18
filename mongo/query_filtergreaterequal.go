package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterGreaterEqual

type query_filterGreaterEqual struct {
	query_filterBase
	selector string
	value    interface{}
}

func (self *query_filterGreaterEqual) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$gte": self.value}}
}

func (self *query_filterGreaterEqual) Selector() string {
	return self.selector
}
