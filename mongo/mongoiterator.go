package mongo

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mgo"
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

func (self *MongoIterator) Next() interface{} {
	if self.iter.Err() != nil {
		return self.iter.Err()
	}
	document := self.collection.NewDocument(self.selectors...)
	if ok := self.iter.Next(document); !ok {
		return nil
	}
	// document has to be initialized again,
	// because mgo zeros the struct while unmarshalling.
	// Newly created slice elements need to be initialized too
	self.collection.InitDocument(document)
	return document
}

func (self *MongoIterator) Err() error {
	return self.iter.Err()
}
