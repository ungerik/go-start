package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterExists

type query_filterExists struct {
	query_filterBase
	selector string
	exists   bool
}

func (self *query_filterExists) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$exists": self.exists}}
}

func (self *query_filterExists) Selector() string {
	return self.selector
}

//func (self *query_filterExists) Exists() bool {
//	return self.exists
//}
