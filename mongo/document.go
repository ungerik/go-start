package mongo

import (
	"github.com/ungerik/go-start/mgo/bson"
	"github.com/ungerik/go-start/model"
)

///////////////////////////////////////////////////////////////////////////////
// Document

type Document interface {
	Init(collection *Collection, embeddingStruct interface{})
	Collection() *Collection
	ObjectId() bson.ObjectId
	SetObjectId(id bson.ObjectId)
	Iterator() model.Iterator
	Ref() Ref
	Save() error
	Remove() error
	// RemoveInvalidRefs sets all invalid mongo.Ref instances to empty,
	// but does not save the document.
	RemoveInvalidRefs() (invalidRefs []InvalidRefData, err error)
}
