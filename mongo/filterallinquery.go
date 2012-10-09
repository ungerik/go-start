package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterAllInQuery

type filterAllInQuery struct {
	filterQueryBase
	selector string
	values   []interface{}
}

func (self *filterAllInQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$all": self.values}}
}

func (self *filterAllInQuery) Selector() string {
	return self.selector
}

//func (self *filterAllInQuery) Values() []interface{} {
//	return self.values
//}
