package model

import (
	"github.com/ungerik/go-start/utils"
)

// CopyFields copies matching struct/array/slice fields
// that implement the Value interface from src to dst.
// Matching is done by comparing the MetaData.Selector()
// of the fields.
// Copying is done by calling SetString() of dst with
// the result from String() of src.
// This way even not exactly matching models can be
// copied by converting to and from strings if this
// is possible without errors.
func CopyFields(dst, src interface{}) (err error) {
	dstFields := map[string]Value{}
	Visit(dst, FieldOnlyVisitor(
		func(field *MetaData) error {
			if dstValue, ok := field.ModelValue(); ok {
				dstFields[field.Selector()] = dstValue
			}
			return nil
		},
	))
	return Visit(src, FieldOnlyVisitor(
		func(field *MetaData) (err error) {
			if srcValue, ok := field.ModelValue(); ok {
				if dstValue, ok := dstFields[field.Selector()]; ok {
					err = dstValue.SetString(srcValue.String())
				}
			}
			return err
		},
	))
}

// CopyFieldsIfNotEmpty copies matching struct/array/slice fields
// that implement the Value interface from src to dst,
// if IsEmpty() of src returns false.
// Matching is done by comparing the MetaData.Selector()
// of the fields.
// Copying is done by calling SetString() of dst with
// the result from String() of src.
// This way even not exactly matching models can be
// copied by converting to and from strings if this
// is possible without errors.
func CopyFieldsIfNotEmpty(dst, src interface{}) (err error) {
	dstFields := map[string]Value{}
	Visit(dst, FieldOnlyVisitor(
		func(field *MetaData) error {
			if dstValue, ok := field.ModelValue(); ok {
				dstFields[field.Selector()] = dstValue
			}
			return nil
		},
	))
	return Visit(src, FieldOnlyVisitor(
		func(field *MetaData) (err error) {
			if srcValue, ok := field.ModelValue(); ok && !srcValue.IsEmpty() {
				if dstValue, ok := dstFields[field.Selector()]; ok {
					err = dstValue.SetString(srcValue.String())
				}
			}
			return err
		},
	))
}

// SetAllSliceLengths sets the length of all slices in document.
func SetAllSliceLengths(document interface{}, length int) {
	Visit(document, VisitorFunc(
		func(data *MetaData) error {
			if data.Kind == SliceKind && data.Value.Len() != length {
				data.Value.Set(utils.SetSliceLengh(data.Value, length))
			}
			return nil
		},
	))
}
