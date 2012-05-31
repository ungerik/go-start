package model

import "fmt"

func NewRequiredError(metaData *MetaData) *RequiredError {
	return &RequiredError{metaData}
}

type RequiredError struct {
	metaData *MetaData
}

func (self RequiredError) Error() string {
	return fmt.Sprintf("Field '%s' is required", self.metaData.Selector())
}
