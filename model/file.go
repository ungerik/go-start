package model

type File struct {
	Name string
	Data []byte
}

func (self *File) String() string {
	return self.Name
}

func (self *File) SetString(str string) error {
	self.Name = str
	return nil
}

func (self *File) IsEmpty() bool {
	return len(self.Data) == 0
}

func (self *File) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *File) Validate(metaData *MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	return nil
}
