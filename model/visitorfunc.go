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

func (self VisitorFunc) BeginIndexedFields(indexedFields *MetaData) error {
	return self(indexedFields)
}

func (self VisitorFunc) IndexedField(field *MetaData) error {
	return self(field)
}

func (self VisitorFunc) EndIndexedFields(indexedFields *MetaData) error {
	return self(indexedFields)
}
