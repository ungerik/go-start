package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterGreater

type query_filterGreater struct {
	query_filterBase
	selector string
	value    interface{}
}

func (self *query_filterGreater) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$gt": self.value}}
}

func (self *query_filterGreater) Selector() string {
	return self.selector
}
