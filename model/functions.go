package model

import (
	"github.com/ungerik/go-start/debug"
	"reflect"
	"strconv"
	"unicode"
)

func init() {
	debug.Nop()
}

func Validate(data interface{}, maxDepth int) []*ValidationError {
	errors := []*ValidationError{}
	WalkStructure(data, maxDepth, func(data interface{}, metaData *MetaData) {
		if validator, ok := data.(Validator); ok {
			errors = append(errors, validator.Validate(metaData)...)
		}
	})
	return errors
}

//func HasParent(value reflect.Value, metaData *MetaData) bool {
//	if value.Kind() == reflect.Struct {
//		for i := range metaData {
//			if metaData[i].ParentStruct.Pointer() == value.Pointer() {
//				return true
//			}
//		}
//	}
//	return false
//}

type WalkStructureCallback func(data interface{}, metaData *MetaData)

// If maxDepth is zero, no limit will be used
func WalkStructure(data interface{}, maxDepth int, callback WalkStructureCallback) {
	if maxDepth == 0 {
		maxDepth = int(^uint(0) >> 1) // Max int
	}
	metaData := &MetaData{}
	walkStructure(reflect.ValueOf(data), metaData, maxDepth, callback)
}

func walkStructure(v reflect.Value, metaData *MetaData, maxDepth int, callback WalkStructureCallback) {
	switch v.Kind() {
	case reflect.Struct:
		if metaData.Depth > maxDepth {
			break
		}
		if v.CanAddr() {
			if _, ok := v.Addr().Interface().(Reference); ok {
				break // Don't go deeper into references
			}
		}
		n := v.NumField()
		for i := 0; i < n; i++ {
			fieldType := v.Type().Field(i)
			// Only walk exported fields
			if unicode.IsUpper(rune(fieldType.Name[0])) {
				m := metaData
				if !fieldType.Anonymous {
					m = &MetaData{
						Parent:       metaData,
						Depth:        metaData.Depth + 1,
						ParentStruct: v,
						Name:         fieldType.Name,
						Index:        -1,
						tag:          fieldType.Tag.Get("gostart"),
					}
				}
				walkStructure(v.Field(i), m, maxDepth, callback)
			}
		}

	case reflect.Slice, reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			m := &MetaData{
				Parent:       metaData,
				Depth:        metaData.Depth + 1,
				ParentStruct: v,
				Name:         strconv.Itoa(i),
				Index:        i,
			}
			walkStructure(v.Index(i), m, maxDepth, callback)
		}
		return

	case reflect.Func, reflect.Map, reflect.Chan, reflect.Invalid:
		return

	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			// TODO: Don't go into cyclic reference when value has itself as parent
			walkStructure(v.Elem(), metaData, maxDepth, callback)
		}
		return
	}

	// v.Addr() to create a pointer type to enable changing the value and
	// casting to struct types whose methods use pointer to struct
	if v.CanAddr() {
		callback(v.Addr().Interface(), metaData)
	}
}
