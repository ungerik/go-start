package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterRefQuery

type filterRefQuery struct {
	filterQueryBase
	selector string
	refs     []Ref
}

func (self *filterRefQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$in": self.refs}}
}

func (self *filterRefQuery) Selector() string {
	return self.selector
}
