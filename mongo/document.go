package mongo

import (
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/ungerik/go-start/model"
)

///////////////////////////////////////////////////////////////////////////////
// Document

type Document interface {
	Init(collection *Collection, embeddingStruct interface{})
	Collection() *Collection
	ObjectId() bson.ObjectId
	SetObjectId(id bson.ObjectId)
	// Iterator returns an iterator that iterates only this document.
	Iterator() model.Iterator
	CreationTime() time.Time
	Ref() Ref
	Save() error
	Delete() error

	// RemoveInvalidRefs sets all invalid mongo.Ref instances to empty,
	// but does not save the document.
	RemoveInvalidRefs() (invalidRefs []InvalidRefData, err error)
}
