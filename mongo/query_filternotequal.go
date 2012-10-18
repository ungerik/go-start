package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterNotEqual

type query_filterNotEqual struct {
	query_filterBase
	selector string
	value    interface{}
}

func (self *query_filterNotEqual) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$ne": self.value}}
}

func (self *query_filterNotEqual) Selector() string {
	return self.selector
}
