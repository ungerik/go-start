package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_filterModulo

type query_filterModulo struct {
	query_filterBase
	selector string
	divisor  interface{}
	result   interface{}
}

func (self *query_filterModulo) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$mod": []interface{}{self.divisor, self.result}}}
}

func (self *query_filterModulo) Selector() string {
	return self.selector
}
