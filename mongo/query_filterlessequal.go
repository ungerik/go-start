package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterLessEqual

type query_filterLessEqual struct {
	query_filterBase
	selector string
	value    interface{}
}

func (self *query_filterLessEqual) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$lte": self.value}}
}

func (self *query_filterLessEqual) Selector() string {
	return self.selector
}
