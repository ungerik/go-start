package mongo

import (
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// query_subDocument

type query_subDocument struct {
	query_filterBase
	selector string
}

func (self *query_subDocument) subDocumentSelector() string {
	return self.selector
}

func (self *query_subDocument) bsonSelector() bson.M {
	return bson.M{self.selector: 1}
}

func (self *query_subDocument) Selector() string {
	return "" // Empty because the self.selector is already returned by subDocumentSelector()
}
