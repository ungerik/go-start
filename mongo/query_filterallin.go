package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterAllIn

type query_filterAllIn struct {
	query_filterBase
	selector string
	values   []interface{}
}

func (self *query_filterAllIn) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$all": self.values}}
}

func (self *query_filterAllIn) Selector() string {
	return self.selector
}

//func (self *query_filterAllIn) Values() []interface{} {
//	return self.values
//}
