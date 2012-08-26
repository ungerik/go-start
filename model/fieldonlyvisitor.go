package model

// FieldOnlyVisitor calls its function for every struct, array and slice field.
type FieldOnlyVisitor func(field *MetaData) error

func (self FieldOnlyVisitor) BeginNamedFields(namedFields *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) NamedField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndNamedFields(namedFields *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) BeginSlice(slice *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) SliceField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndSlice(slice *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) BeginArray(array *MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) ArrayField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndArray(array *MetaData) error {
	return nil
}
