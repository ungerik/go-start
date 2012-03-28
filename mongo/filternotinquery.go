package mongo

import (
	"launchpad.net/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterNotInQuery

type filterNotInQuery struct {
	filterQueryBase
	selector string
	values   []interface{}
}

func (self *filterNotInQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$nin": self.values}}
}

func (self *filterNotInQuery) Selector() string {
	return self.selector
}
