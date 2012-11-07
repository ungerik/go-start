package model

type Value interface {
	String() string
	// SetString returns only error from converting str to the
	// underlying value type.
	// It does not return validation errors of the converted value.
	SetString(str string) (strconvErr error)
	IsEmpty() bool
	Required(metaData *MetaData) bool
	Validator
}
