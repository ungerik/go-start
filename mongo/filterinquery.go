package mongo

import (
	"launchpad.net/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterInQuery

type filterInQuery struct {
	filterQueryBase
	selector string
	values   []interface{}
}

func (self *filterInQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$in": self.values}}
}

func (self *filterInQuery) Selector() string {
	return self.selector
}
