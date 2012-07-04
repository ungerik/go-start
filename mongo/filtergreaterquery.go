package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterGreaterQuery

type filterGreaterQuery struct {
	filterQueryBase
	selector string
	value    interface{}
}

func (self *filterGreaterQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$gt": self.value}}
}

func (self *filterGreaterQuery) Selector() string {
	return self.selector
}
