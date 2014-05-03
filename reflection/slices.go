package reflection

import (
	"fmt"
	"reflect"
)

// AppendDefaultSliceElement appends a field to slice with the
// value returned by GetDefaultValue for the former last slice
// element, or the zero value of the type if the slice was empty.
func AppendDefaultSliceElement(slice reflect.Value) reflect.Value {
	if slice.Kind() != reflect.Slice {
		panic(fmt.Errorf("Expected slice type, got %T", slice.Interface()))
	}
	var newField reflect.Value
	if slice.Len() > 0 {
		newField = GetDefaultValue(slice.Index(slice.Len() - 1))
	} else {
		newField = reflect.Zero(slice.Type().Elem())
	}
	return reflect.Append(slice, newField)
}

// SetSliceLengh sets the length of a slice by sub-slicing a slice that's too long,
// or appending empty fields with the result of AppendDefaultSliceElement.
func SetSliceLengh(slice reflect.Value, length int) reflect.Value {
	if length > slice.Len() {
		for i := slice.Len(); i < length; i++ {
			slice = AppendDefaultSliceElement(slice)
		}
	} else if length < slice.Len() {
		slice = slice.Slice(0, length)
	}
	return slice
}

// DeleteDefaultSliceElementsVal deletes slice elements where IsDefaultValue
// returns true.
func DeleteDefaultSliceElementsVal(slice reflect.Value) reflect.Value {
	if slice.Kind() != reflect.Slice {
		panic(fmt.Errorf("Expected slice type, got %T", slice.Interface()))
	}
	for i := slice.Len() - 1; i >= 0; i-- {
		if IsDefaultValue(slice.Index(i)) {
			fmt.Println("Found default", i)
			before := slice.Slice(0, i)
			if i == slice.Len()-1 {
				slice = before
			} else {
				after := slice.Slice(i+1, slice.Len())
				slice = reflect.AppendSlice(before, after)
			}
			i--
		}
	}
	return slice
}

// DeleteDefaultSliceElements deletes slice elements where IsDefaultValue
// returns true.
func DeleteDefaultSliceElements(slice interface{}) interface{} {
	return DeleteDefaultSliceElementsVal(reflect.ValueOf(slice)).Interface()
}
