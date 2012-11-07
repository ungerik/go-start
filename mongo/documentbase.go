package mongo

import (
	"reflect"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
)

///////////////////////////////////////////////////////////////////////////////
// DocumentBase

type DocumentBase struct {
	ID              bson.ObjectId `bson:"_id,omitempty" gostart:"-"`
	collection      *Collection   `gostart:"-"`
	embeddingStruct interface{}   `gostart:"-"`
}

func (self *DocumentBase) Init(collection *Collection, embeddingStructPtr interface{}) {
	if collection == nil {
		panic("mongo.DocumentBase.Init() called with collection == nil")
	}
	if embeddingStructPtr == nil {
		panic("mongo.DocumentBase.Init() called with embeddingStructPtr == nil")
	}
	if reflect.ValueOf(embeddingStructPtr).Kind() != reflect.Ptr {
		panic("mongo.DocumentBase.Init() embeddingStructPtr must be pointer, but is " + reflect.ValueOf(embeddingStructPtr).Elem().Type().String())
	}

	self.collection = collection
	self.embeddingStruct = embeddingStructPtr

	InitRefs(embeddingStructPtr)
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

func (self *DocumentBase) CreationTime() (t time.Time) {
	if !self.ID.Valid() {
		return t
	}
	return self.ID.Time()
}

func (self *DocumentBase) Ref() Ref {
	return Ref{self.ID, self.collection.Name}
}

func (self *DocumentBase) Save() (err error) {
	if self.embeddingStruct == nil {
		return errs.Format("Can't save uninitialized mongo.Document. embeddingStruct is nil.")
	}
	// if self.ID.Valid() {
	// 	err = self.collection.collection.UpdateId(self.ID, self.embeddingStruct)
	// } else {
	// 	self.ID = bson.NewObjectId()
	// 	err = self.collection.collection.Insert(self.embeddingStruct)
	// }
	// if err == nil {
	// 	n, err := self.collection.FilterID(self.ID).Count()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if n != 1 {
	// 		// panic("Something went wrong with saving")
	// 		debug.Print("Something went wrong with saving")
	// 	}
	// }
	info, err := self.collection.collection.UpsertId(self.ID, self.embeddingStruct)
	if err != nil {
		return err
	}
	if info.UpsertedId != nil {
		self.ID = info.UpsertedId.(bson.ObjectId)
	} else if info.Updated != 1 {
		debug.Nop()
		return errs.Format("Something went wrong with saving")
	}
	return nil
}

// func (self *DocumentBase) Save() error {
// 	if self.embeddingStruct == nil {
// 		return errs.Format("Can't save uninitialized mongo.Document. embeddingStruct is nil.")
// 	}
// 	// if !self.ID.Valid() {
// 	// 	self.ID = bson.NewObjectId()
// 	// }
// 	debug.Dump(self.ID)
// 	info, err := self.collection.collection.UpsertId(self.ID, self.embeddingStruct)
// 	debug.Dump(self.ID)
// 	debug.Dump(info)
// 	if err != nil {
// 		return err
// 	}
// 	// if info.UpsertedId != nil && info.UpsertedId.(bson.ObjectId) != self.ID {
// 	// 	panic("Something is wrong with Upsert")
// 	// }
// 	return nil
// }

func (self *DocumentBase) Delete() error {
	return self.collection.DeleteWithID(self.ID)
}

func (self *DocumentBase) RemoveInvalidRefs() (invalidRefs []InvalidRefData, err error) {
	return RemoveInvalidRefs(self.embeddingStruct)
}
