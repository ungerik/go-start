package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterLessQuery

type filterLessQuery struct {
	filterQueryBase
	selector string
	value    interface{}
}

func (self *filterLessQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$lt": self.value}}
}

func (self *filterLessQuery) Selector() string {
	return self.selector
}
