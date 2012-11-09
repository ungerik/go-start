package media

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

func NewBlobRef(blob *Blob) *BlobRef {
	blobRef := new(BlobRef)
	blobRef.Set(blob)
	return blobRef
}

type BlobRef string

func (self *BlobRef) String() string {
	return string(*self)
}

func (self *BlobRef) SetString(str string) error {
	*self = BlobRef(str)
	return nil
}

// Blob loads the referenced blob, or returns nil if the reference is empty.
func (self *BlobRef) Get() (*Blob, error) {
	if self.IsEmpty() {
		return nil, nil
	}
	return Config.Backend.LoadBlob(self.String())
}

func (self *BlobRef) TryGet() (*Blob, bool, error) {
	if self.IsEmpty() {
		return nil, false, nil
	}
	blob, err := Config.Backend.LoadBlob(self.String())
	if _, notFound := err.(ErrNotFound); notFound {
		return nil, false, nil
	}
	return blob, err == nil, err
}

// SetBlob sets the ID of blob, or an empty reference if blob is nil.
func (self *BlobRef) Set(blob *Blob) {
	if blob != nil {
		self.SetString(blob.ID.String())
	} else {
		*self = ""
	}
}

func (self *BlobRef) IsEmpty() bool {
	return *self == ""
}

func (self *BlobRef) Required(metaData *model.MetaData) bool {
	return metaData.BoolAttrib(model.StructTagKey, "required")
}

func (self *BlobRef) Validate(metaData *model.MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return model.NewRequiredError(metaData)
	}
	return nil
}

func (self *BlobRef) FileURL() (view.URL, error) {
	blob, err := self.Get()
	if err != nil {
		return nil, err
	}
	return blob.FileURL(), nil
}

func (self *BlobRef) FileLink(class string) (*view.Link, error) {
	blob, err := self.Get()
	if err != nil {
		return nil, err
	}
	return blob.FileLink(class), nil
}
