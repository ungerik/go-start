package mongo

import (
	"github.com/ungerik/go-start/model"
	"labix.org/v2/mgo"
)

///////////////////////////////////////////////////////////////////////////////
// MongoIterator

func NewMongoIterator(query Query) model.Iterator {
	mgoQuery, err := query.mongoQuery()
	if err != nil {
		return model.NewErrorOnlyIterator(err)
	}
	return &MongoIterator{collection: query.Collection(), iter: mgoQuery.Iter()}
	// mgoIter := mgoQuery.Iter()
	// collection, selectors := collectionAndSubDocumentSelectors(query)
	// return &MongoIterator{collection: collection, selectors: selectors, iter: mgoIter}
}

type MongoIterator struct {
	collection *Collection
	// selectors  []string
	iter *mgo.Iter
}

func (self *MongoIterator) Next(resultRef interface{}) bool {
	if self.iter.Err() != nil {
		return false
	}
	if !self.iter.Next(resultRef) {
		return false
	}
	// resultRef has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.collection.InitDocument(documentFromResultPtr(resultRef))
	return true
}

func (self *MongoIterator) Err() error {
	return self.iter.Err()
}
