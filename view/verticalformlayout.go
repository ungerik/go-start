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
}

func (self *VerticalFormLayout) BeginFormContent(form *Form) View {
	return nil
}

func (self *VerticalFormLayout) EndFormContent(form *Form) View {
	return nil
}

func (self *VerticalFormLayout) BeginStruct(form *Form, strct *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) StructField(form *Form, field *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) EndStruct(form *Form, strct *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) BeginArray(form *Form, array *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) ArrayField(form *Form, field *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) EndArray(form *Form, array *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) BeginSlice(form *Form, slice *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) SliceField(form *Form, field *model.MetaData) View {
	return nil
}

func (self *VerticalFormLayout) EndSlice(form *Form, slice *model.MetaData) View {
	return nil
}
