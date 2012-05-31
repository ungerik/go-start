package model

type Validator interface {
	// Validate returns an error if metaData is not valid.
	// In case of multiple errors errs.MultipleErrors is returned.
	Validate(metaData *MetaData) error
}
