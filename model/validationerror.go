package model

import (
	"fmt"
	"github.com/ungerik/go-start/errs"
)

var NoValidationErrors []*ValidationError = []*ValidationError{}

type ValidationError struct {
	WrappedError error
	MetaData     MetaData
}

func (self *ValidationError) Error() string {
	return fmt.Sprintf("ValidationError %s: %s", self.MetaData.Selector(), self.WrappedError)
}

func NewValidationErrors(error error, metaData MetaData) []*ValidationError {
	return []*ValidationError{&ValidationError{error, metaData}}
}

func NewRequiredValidationError(metaData MetaData) *ValidationError {
	return &ValidationError{
		errs.Format("Field '%s' is required", metaData.Selector()),
		metaData,
	}
}
