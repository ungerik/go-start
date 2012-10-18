package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterNotIn

type query_filterNotIn struct {
	query_filterBase
	selector string
	values   []interface{}
}

func (self *query_filterNotIn) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$nin": self.values}}
}

func (self *query_filterNotIn) Selector() string {
	return self.selector
}
