package model

type Value interface {
	//Init(metaData *MetaData)
	String() string
	SetString(str string) error
	IsEmpty() bool
	Required(metaData *MetaData) bool
	Validator
}
