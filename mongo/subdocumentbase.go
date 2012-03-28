package mongo

import (
	"launchpad.net/mgo/bson"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
)

///////////////////////////////////////////////////////////////////////////////
// SubDocumentBase

type SubDocumentBase struct {
	rootDocumentID  bson.ObjectId
	collection      *Collection
	selector        string
	embeddingStruct interface{}
}

func (self *SubDocumentBase) Init(collection *Collection, selector string, embeddingStruct interface{}) {
	if collection == nil {
		panic("mongo.SubDocumentBase.Init() called with nil collection")
	}
	if selector == "" {
		panic("mongo.SubDocumentBase.Init() called with empty selector")
	}
	if embeddingStruct == nil {
		panic("mongo.SubDocumentBase.Init() called with nil embeddingStruct")
	}

	self.collection = collection
	self.selector = selector
	self.embeddingStruct = embeddingStruct

	InitRefs(embeddingStruct)
}

func (self *SubDocumentBase) Collection() *Collection {
	return self.collection
}

// Implements DocumentInterface
func (self *SubDocumentBase) RootDocumentObjectId() bson.ObjectId {
	return self.rootDocumentID
}

func (self *SubDocumentBase) RootDocumentSetObjectId(id bson.ObjectId) {
	self.rootDocumentID = id
}

func (self *SubDocumentBase) Iterator() model.Iterator {
	return model.NewObjectIterator(self.embeddingStruct)
}

func (self *SubDocumentBase) Save() error {
	if self.embeddingStruct == nil {
		return errs.Format("Can't save uninitialized mongo.SubDocument. embeddingStruct is nil.")
	}
	if !self.rootDocumentID.Valid() {
		return errs.Format("Can't save mongo.SubDocument with invalid RootDocumentObjectId.")
	}

	return self.collection.Update(self.rootDocumentID, bson.M{self.selector: self.embeddingStruct})
}

func (self *SubDocumentBase) RemoveRootDocument() error {
	return self.collection.Remove(self.rootDocumentID)
}
