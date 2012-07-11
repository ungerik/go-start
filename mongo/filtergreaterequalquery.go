package mongo

import (
	"github.com/ungerik/go-start/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterGreaterEqualQuery

type filterGreaterEqualQuery struct {
	filterQueryBase
	selector string
	value    interface{}
}

func (self *filterGreaterEqualQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$gte": self.value}}
}

func (self *filterGreaterEqualQuery) Selector() string {
	return self.selector
}
