package model

// FieldOnlyVisitor calls its function for every struct, array and slice field.
type FieldOnlyVisitor func(field *MetaData) error

func (self FieldOnlyVisitor) BeginNamedFields(*MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) NamedField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndNamedFields(*MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) BeginIndexedFields(*MetaData) error {
	return nil
}

func (self FieldOnlyVisitor) IndexedField(field *MetaData) error {
	return self(field)
}

func (self FieldOnlyVisitor) EndIndexedFields(*MetaData) error {
	return nil
}
