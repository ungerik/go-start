package mongo

import (
	"launchpad.net/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterArraySizeQuery

type filterArraySizeQuery struct {
	filterQueryBase
	selector string
	size     int
}

func (self *filterArraySizeQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$size": self.size}}
}

func (self *filterArraySizeQuery) Selector() string {
	return self.selector
}

//func (self *filterArraySizeQuery) Size() int {
//	return self.size
//}

