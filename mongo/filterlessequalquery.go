package mongo

import (
	"github.com/ungerik/go-start/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterLessEqualQuery

type filterLessEqualQuery struct {
	filterQueryBase
	selector string
	value    interface{}
}

func (self *filterLessEqualQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$lte": self.value}}
}

func (self *filterLessEqualQuery) Selector() string {
	return self.selector
}
