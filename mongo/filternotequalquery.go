package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterNotEqualQuery

type filterNotEqualQuery struct {
	filterQueryBase
	selector string
	value    interface{}
}

func (self *filterNotEqualQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$ne": self.value}}
}

func (self *filterNotEqualQuery) Selector() string {
	return self.selector
}
