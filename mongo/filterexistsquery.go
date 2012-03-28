package mongo

import (
	"launchpad.net/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterExistsQuery

type filterExistsQuery struct {
	filterQueryBase
	selector string
	exists   bool
}

func (self *filterExistsQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$exists": self.exists}}
}

func (self *filterExistsQuery) Selector() string {
	return self.selector
}

//func (self *filterExistsQuery) Exists() bool {
//	return self.exists
//}

