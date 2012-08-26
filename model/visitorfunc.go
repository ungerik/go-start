package model

// VisitorFunc calls its function for every Visitor method call,
// thus mapping all Visitor methods on a single function.
type VisitorFunc func(data *MetaData) error

func (self VisitorFunc) BeginNamedFields(namedFields *MetaData) error {
	return self(namedFields)
}

func (self VisitorFunc) NamedField(field *MetaData) error {
	return self(field)
}

func (self VisitorFunc) EndNamedFields(namedFields *MetaData) error {
	return self(namedFields)
}

func (self VisitorFunc) BeginSlice(slice *MetaData) error {
	return self(slice)
}

func (self VisitorFunc) SliceField(field *MetaData) error {
	return self(field)
}

func (self VisitorFunc) EndSlice(slice *MetaData) error {
	return self(slice)
}

func (self VisitorFunc) BeginArray(array *MetaData) error {
	return self(array)
}

func (self VisitorFunc) ArrayField(field *MetaData) error {
	return self(field)
}

func (self VisitorFunc) EndArray(array *MetaData) error {
	return self(array)
}
