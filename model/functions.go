package model

import (
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/errs"
	// "reflect"
	// "strconv"
	// "unicode"
)

func init() {
	debug.Nop()
}

func Validate(data interface{}, maxDepth int) error {
	var errors []error
	VisitMaxDepth(data, maxDepth, VisitorFunc(
		func(data *MetaData) error {
			if validator, ok := data.Value.Addr().Interface().(Validator); ok {
				errors = append(errors, validator.Validate(data))
			}
			return nil
		},
	))
	return errs.Errors(errors...)
}

/* 

todo remove

// WalkStructure walks recursively over all fields of data and
// reports them depth first via callback.
// If maxDepth is zero, no limit will be used
func WalkStructure(data interface{}, maxDepth int, callback func(data *MetaData)) {
	if maxDepth == 0 {
		maxDepth = int(^uint(0) >> 1) // Max int
	}
	walkStructure(&MetaData{Value: reflect.ValueOf(data)}, maxDepth, callback)
}

func walkStructure(data *MetaData, maxDepth int, callback func(data *MetaData)) {
	v := data.Value
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			child := &MetaData{
				Value:  v.Index(i),
				Parent: data,
				Depth:  data.Depth + 1,
				Name:   strconv.Itoa(i),
				Index:  i,
			}
			walkStructure(child, maxDepth, callback)
		}
		return

	case reflect.Func, reflect.Map, reflect.Chan, reflect.Invalid:
		return

	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			data.Value = v.Elem()
			// TODO: Don't go into cyclic reference when value has itself as parent
			walkStructure(data, maxDepth, callback)
		}
		return

	case reflect.Struct:
		if data.Depth > maxDepth {
			break
		}
		if v.CanAddr() {
			if _, ok := v.Addr().Interface().(Reference); ok {
				break // Don't go deeper into references
			}
		}
		for i := 0; i < v.NumField(); i++ {
			fieldType := v.Type().Field(i)
			// Only walk exported fields
			if unicode.IsUpper(rune(fieldType.Name[0])) {
				child := data
				if !fieldType.Anonymous {
					child = &MetaData{
						Value:  v.Field(i),
						Parent: data,
						Depth:  data.Depth + 1,
						Name:   fieldType.Name,
						Index:  -1,
						tag:    fieldType.Tag.Get("gostart"),
					}
				}
				walkStructure(child, maxDepth, callback)
			}
		}
	}

	if !v.CanAddr() {
		// will this ever be the case???
		return
	}

	callback(data)
}

*/
