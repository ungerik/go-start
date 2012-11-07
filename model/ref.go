package model

import (
	"labix.org/v2/mgo/bson"
)

type Ref string

func (self *Ref) String() string {
	return string(*self)
}

func (self *Ref) SetString(str string) (strconvErr error) {
	*self = Ref(str)
	return nil
}

func (self *Ref) IsEmpty() bool {
	return *self == ""
}

func (self *Ref) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *Ref) StringID() string {
	return self.String()
}

// Implements bson.Setter for mongo.Ref compatibility.
// Yes this bson dependency is dirty, but necessary
// for transition form mongo.Ref to model.Ref.
func (self *Ref) SetBSON(raw bson.Raw) error {
	var objectId *bson.ObjectId
	if raw.Unmarshal(&objectId) == nil {
		if objectId != nil {
			*self = Ref(objectId.Hex())
		} else {
			*self = ""
		}
		return nil
	}
	return raw.Unmarshal(self)
}
