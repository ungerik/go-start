package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterIn

type query_filterIn struct {
	query_filterBase
	selector string
	values   []interface{}
}

func (self *query_filterIn) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$in": self.values}}
}

func (self *query_filterIn) Selector() string {
	return self.selector
}
