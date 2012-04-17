package model

type Validator interface {
	Validate(metaData *MetaData) []*ValidationError	
}

