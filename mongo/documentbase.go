package mongo

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
	"labix.org/v2/mgo/bson"
)

///////////////////////////////////////////////////////////////////////////////
// DocumentBase

type DocumentBase struct {
	ID              bson.ObjectId `bson:"_id,omitempty" gostart:"-"`
	collection      *Collection   `gostart:"-"`
	embeddingStruct interface{}   `gostart:"-"`
}

func (self *DocumentBase) Init(collection *Collection, embeddingStruct interface{}) {
	if collection == nil {
		panic("mongo.DocumentBase.Init() called with collection == nil")
	}
	if embeddingStruct == nil {
		panic("mongo.DocumentBase.Init() called with embeddingStruct == nil")
	}

	self.collection = collection
	self.embeddingStruct = embeddingStruct

	InitRefs(embeddingStruct)
}

func (self *DocumentBase) Collection() *Collection {
	if self == nil {
		return nil
	}
	return self.collection
}

func (self *DocumentBase) ObjectId() bson.ObjectId {
	if self == nil {
		return ""
	}
	return self.ID
}

func (self *DocumentBase) SetObjectId(id bson.ObjectId) {
	self.ID = id
}

func (self *DocumentBase) Iterator() model.Iterator {
	return model.NewSliceIterator(self.embeddingStruct)
}

func (self *DocumentBase) Ref() Ref {
	return Ref{self.ID, self.collection.Name}
}

func (self *DocumentBase) Save() error {
	if self.embeddingStruct == nil {
		return errs.Format("Can't save uninitialized mongo.Document. embeddingStruct is nil.")
	}
	if !self.ID.Valid() {
		self.ID = bson.NewObjectId()
	}
	info, err := self.collection.collection.UpsertId(self.ID, self.embeddingStruct)
	if err != nil {
		return err
	}
	if info.UpsertedId != nil && info.UpsertedId.(bson.ObjectId) != self.ID {
		panic("Something is wrong with Upsert")
	}
	return nil
}

func (self *DocumentBase) Delete() error {
	return self.collection.Delete(self.ID)
}

func (self *DocumentBase) RemoveInvalidRefs() (invalidRefs []InvalidRefData, err error) {
	return RemoveInvalidRefs(self.embeddingStruct)
}
