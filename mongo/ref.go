package mongo

import (
	"fmt"

	"labix.org/v2/mgo/bson"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/model"
)

/*
Ref holds a mongo.ID and the name of the collection of the referenced document.
Only the ID will be saved in MongoDB, the collection name is for validation
and convenience functions only.
*/
type Ref struct {
	ID             bson.ObjectId `gostart:"-"`
	CollectionName string        `gostart:"-"`
}

func (self *Ref) StringID() string {
	return self.ID.Hex()
}

func (self *Ref) SetStringID(str string) error {
	switch len(str) {
	case 0, 12:
		self.ID = bson.ObjectId(str)
	case 24:
		self.ID = bson.ObjectIdHex(str)
	default:
		return errs.Format("Invalid string for bson.ObjectId: '%s'", str)
	}
	return nil
}

func (self *Ref) String() string {
	return fmt.Sprintf("mongo.Ref of ID %s in collection %s", self.ID.Hex(), self.CollectionName)
}

func (self *Ref) DocumentLabel() (string, error) {
	return self.Collection().DocumentLabel(self.ID)
}

// Returns nil and no error if the reference is empty
func (self *Ref) Get(resultRef interface{}) error {
	return self.Collection().DocumentWithID(self.ID, resultRef)
}

// Returns an error if the reference is empty
func (self *Ref) TryGet(resultRef interface{}) (found bool, err error) {
	return self.Collection().TryDocumentWithID(self.ID, resultRef)
}

// nil is valid and sets the reference to empty
func (self *Ref) Set(document Document) {
	if document == nil {
		self.ID = ""
		return
	}
	if document.Collection() == nil {
		panic("model.Document.Collection() == nil")
	}
	if document.Collection().Name != self.CollectionName {
		panic(fmt.Sprintf("Can't set document from different collection. Expected collection '%s', got '%s'", self.CollectionName, document.Collection().Name))
	}
	self.ID = document.ObjectId()
}

func (self *Ref) IsEmpty() bool {
	return self.ID == ""
}

// Implements bson.Getter
func (self Ref) GetBSON() (interface{}, error) {
	if !self.ID.Valid() {
		return nil, nil
	}
	return self.ID, nil
}

// Implements bson.Setter
func (self *Ref) SetBSON(raw bson.Raw) error {
	var id *bson.ObjectId
	err := raw.Unmarshal(&id)
	if err != nil {
		return err
	}
	if id == nil {
		self.ID = ""
	} else {
		self.ID = *id
	}
	return nil
}

func (self *Ref) Validate(metaData *model.MetaData) error {
	if self.CollectionName == "" {
		return errs.Format("Missing CollectionName")
	}
	length := len(self.ID)
	if length != 0 && length != 12 {
		return errs.Format("Invalid ObjectId length %d", length)
	}
	if self.Required(metaData) && self.IsEmpty() {
		return model.NewRequiredError(metaData)
	}
	if ok, err := self.CheckID(); !ok {
		if err != nil {
			return err
		}
		return errs.Format("Invalid mongo.Ref. No document with ID %s in collection %s", self.ID, self.CollectionName)
	}
	return nil
}

func (self *Ref) CheckID() (bool, error) {
	if self.IsEmpty() {
		return true, nil
	}
	count, err := self.Collection().FilterReferenced(*self).Count()
	return count == 1, err
}

func (self *Ref) Required(metaData *model.MetaData) bool {
	return metaData.BoolAttrib(model.StructTagKey, "required")
}

func (self *Ref) Collection() *Collection {
	if self.CollectionName == "" {
		panic(errs.Format("Missing collection name. Did you call mongo.Document.Init()?"))
	}
	collection, ok := Collections[self.CollectionName]
	if !ok {
		panic(errs.Format("Collection '" + self.CollectionName + "' not registered"))
	}
	return collection
}

func (self *Ref) References(doc Document) (ok bool, err error) {
	collection := doc.Collection()
	if collection == nil {
		return false, errs.Format("Document is not initialized")
	}
	if collection.Name != self.CollectionName {
		return false, errs.Format("mongo.Ref to collection '%s' can't reference document of collection '%s'", self.CollectionName, collection.Name)
	}
	id := doc.ObjectId()
	if !id.Valid() {
		return false, errs.Format("Can't reference document with empty ID")
	}
	return self.ID == id, nil
}
