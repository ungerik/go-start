package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterEqual

type query_filterEqual struct {
	query_filterBase
	selector string
	value    interface{}
}

func (self *query_filterEqual) bsonSelector() bson.M {
	return bson.M{self.selector: self.value}
}

func (self *query_filterEqual) Selector() string {
	// Don't return special selectors that start with $
	if self.selector != "" && self.selector[0] == '$' {
		return ""
	}
	return self.selector
}

//func (self *query_filterEqual) Value() interface{} {
//	return self.value
//}
