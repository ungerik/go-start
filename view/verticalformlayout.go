package view

import (
	"fmt"
	"github.com/ungerik/go-start/model"
	"reflect"
)

/*
VerticalFormLayout.

CSS needed for VerticalFormLayout:

	form label:after {
		content: ":";
	}

	form input[type=checkbox] + label:after {
		content: "";
	}

Additional CSS for labels above input fields (except checkboxes):

	form label {
		display: block;
	}

	form input[type=checkbox] + label {
		display: inline;
	}

DIV classes for coloring:

	form .required {}
	form .error {}
	form .success {}

*/
type VerticalFormLayout struct {
	form *Form
}

func (self *VerticalFormLayout) BeginStruct(strct *model.MetaData) {

}

func (self *VerticalFormLayout) StructField(field *model.MetaData) {

}

func (self *VerticalFormLayout) EndStruct(strct *model.MetaData) {

}

func (self *VerticalFormLayout) BeginSlice(slice *model.MetaData) {

}

func (self *VerticalFormLayout) SliceField(field *model.MetaData) {

}

func (self *VerticalFormLayout) EndSlice(slice *model.MetaData) {

}

func (self *VerticalFormLayout) BeginArray(array *model.MetaData) {

}

func (self *VerticalFormLayout) ArrayField(field *model.MetaData) {

}

func (self *VerticalFormLayout) EndArray(array *model.MetaData) {

}
