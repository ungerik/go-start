package mongo

import (
	"github.com/ungerik/go-start/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterQuery

type filterQuery struct {
	filterQueryBase
	selector string
	value    interface{}
}

func (self *filterQuery) bsonSelector() bson.M {
	return bson.M{self.selector: self.value}
}

func (self *filterQuery) Selector() string {
	// Don't return special selectors that start with $
	if self.selector != "" && self.selector[0] == '$' {
		return ""
	}
	return self.selector
}

//func (self *filterQuery) Value() interface{} {
//	return self.value
//}
