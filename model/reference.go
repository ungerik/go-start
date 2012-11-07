package model

type Reference interface {
	Value	
	StringID() string
}
