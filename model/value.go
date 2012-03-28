package model

type Value interface {
	String() string
	SetString(str string) error
	FixValue(metaData MetaData)
	Validator
}
