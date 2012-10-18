package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterRef

type query_filterRef struct {
	query_filterBase
	selector string
	refs     []Ref
}

func (self *query_filterRef) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$in": self.refs}}
}

func (self *query_filterRef) Selector() string {
	return self.selector
}
