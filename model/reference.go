package model

type Reference interface {
	Value	
	Reference()
}
