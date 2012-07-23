package model

type Validator interface {
	// Validate returns an error if metaData is not valid.
	// In case of multiple errors errs.ErrSlice is returned.
	Validate(metaData *MetaData) error
}
