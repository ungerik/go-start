package mongo

import (
	"github.com/ungerik/go-start/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// filterModuloQuery

type filterModuloQuery struct {
	filterQueryBase
	selector string
	divisor  interface{}
	result   interface{}
}

func (self *filterModuloQuery) bsonSelector() bson.M {
	return bson.M{self.selector: bson.M{"$mod": []interface{}{self.divisor, self.result}}}
}

func (self *filterModuloQuery) Selector() string {
	return self.selector
}
