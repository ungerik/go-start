package mongo

import (
	"github.com/ungerik/go-start/model"
	"labix.org/v2/mgo"
)

///////////////////////////////////////////////////////////////////////////////
// MongoIterator

func newIterator(query Query) model.Iterator {
	mgoQuery, err := query.mongoQuery()
	if err != nil {
		return model.NewErrorOnlyIterator(err)
	}
	mgoIter := mgoQuery.Iter()
	collection, selectors := collectionAndSubDocumentSelectors(query)
	return &MongoIterator{collection: collection, selectors: selectors, iter: mgoIter}
}

type MongoIterator struct {
	collection *Collection
	selectors  []string
	iter       *mgo.Iter
}

func (self *MongoIterator) Next(resultPtr interface{}) bool {
	if self.iter.Err() != nil {
		return false
	}
	if !self.iter.Next(resultPtr) {
		return false
	}
	// resultPtr has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.collection.InitDocument(resultPtr.(Document))
	return true
}

func (self *MongoIterator) Err() error {
	return self.iter.Err()
}
