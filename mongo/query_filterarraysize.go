package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterArraySize

type query_filterArraySize struct {
	query_filterBase
	selector string
	size     int
}

func (self *query_filterArraySize) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$size": self.size}}
}

func (self *query_filterArraySize) Selector() string {
	return self.selector
}

//func (self *query_filterArraySize) Size() int {
//	return self.size
//}
