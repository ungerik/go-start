package model

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/utils"
	"reflect"
)

func Validate(data interface{}, maxDepth int) error {
	var errors []error
	VisitMaxDepth(data, maxDepth, VisitorFunc(
		func(data *MetaData) error {
			if validator, ok := data.ModelValidator(); ok {
				errors = append(errors, validator.Validate(data))
			}
			return nil
		},
	))
	return errs.Errors(errors...)
}

func AppendEmptySliceEnds(strct interface{}) {
	utils.VisitStruct(strct, utils.ModifySliceStructVisitor(
		func(depth int, v reflect.Value) (reflect.Value, error) {
			field := reflect.Zero(v.Type().Elem())
			return reflect.Append(v, field), nil
		},
	))
}

func RemoveEmptySliceEnds(strct interface{}) {
	utils.VisitStruct(strct, utils.ModifySliceStructVisitor(
		func(depth int, v reflect.Value) (reflect.Value, error) {
			for v.Len() > 0 {
				last := v.Index(v.Len() - 1)
				if value, ok := last.Addr().Interface().(Value); ok {
					if !value.IsEmpty() {
						return v, nil
					}
					v = v.Slice(0, v.Len()-1)
				} else if last.Kind() == reflect.Struct {
					for i := 0; i < last.NumField(); i++ {
						if value, ok := v.Field(i).Addr().Interface().(Value); ok {
							if !value.IsEmpty() {
								return v, nil
							}
						}
						// ignore struct fields that that don't implement Value
					}
					v = v.Slice(0, v.Len()-1)
				}
			}
			return v, nil
		},
	))
}
