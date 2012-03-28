package mongo

import (
	"launchpad.net/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// subDocumentQuery

type subDocumentQuery struct {
	filterQueryBase
	selector string
}

func (self *subDocumentQuery) subDocumentSelector() string {
	return self.selector
}

func (self *subDocumentQuery) bsonSelector() bson.M {
	return bson.M{self.selector: 1}
}

func (self *subDocumentQuery) Selector() string {
	return "" // Empty because the self.selector is already returned by subDocumentSelector()
}
