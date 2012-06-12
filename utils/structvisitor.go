package utils

import (
	"reflect"
)

////////////////////////////////////////////////////////////////////////////////
// StructVisitor

type StructVisitor interface {
	BeginStruct(depth int, v reflect.Value) error
	StructField(depth int, v reflect.Value, f reflect.StructField, index int) error
	EndStruct(depth int, v reflect.Value) error

	// PreModifySlice is special: it can change the slice length and
	// returns the altered or unaltered slice as reflect.Value.
	PreModifySlice(depth int, v reflect.Value) (reflect.Value, error)
	BeginSlice(depth int, v reflect.Value) error
	SliceField(depth int, v reflect.Value, index int) error
	EndSlice(depth int, v reflect.Value) error
	// PostModifySlice is special: it can change the slice length and
	// returns the altered or unaltered slice as reflect.Value.
	PostModifySlice(depth int, v reflect.Value) (reflect.Value, error)

	BeginArray(depth int, v reflect.Value) error
	ArrayField(depth int, v reflect.Value, index int) error
	EndArray(depth int, v reflect.Value) error
}

/*
VisitStruct visits recursively all exported fields of a struct
and reports them via StructVisitor methods.
If a StructVisitor method returns an error, the visitation is aborted
and the error returned as result.
Pointers and interfaces are dereferenced silently until a non nil value
is found.
Structs that are embedded anonymously are inlined so that their fields
are reported as fields of the embedding struct at the same depth.
Anonymous struct fields that are not structs themselves are omitted.
Struct fields with the tag gostart:"-" are ignored.
*/
func VisitStruct(strct interface{}, visitor StructVisitor) error {
	return VisitStructDepth(strct, visitor, -1)
}

/*
VisitStructDepth is identical to VisitStruct except that its recursive
depth is limited to maxDepth with the first depth level being zero.
If maxDepth is -1, then the recursive depth is unlimited (VisitStruct).
*/
func VisitStructDepth(strct interface{}, visitor StructVisitor, maxDepth int) error {
	return visitStructRecursive(reflect.ValueOf(strct), visitor, maxDepth, 0)
}

func visitAnonymousStructFieldRecursive(visitor StructVisitor, v reflect.Value, index *int, depth int) (err error) {
	if v.Kind() == reflect.Struct {
		t := v.Type()
		n := t.NumField()
		for i := 0; i < n; i++ {
			f := t.Field(i)
			if f.PkgPath == "" && f.Tag.Get("gostart") != "-" { // Only exported fields
				if vi, ok := DereferenceValue(v.Field(i)); ok {
					if f.Anonymous {
						err = visitAnonymousStructFieldRecursive(visitor, vi, index, depth)
					} else {
						err = visitor.StructField(depth, vi, f, *index)
						*index++
					}
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func visitStructRecursive(v reflect.Value, visitor StructVisitor, maxDepth, depth int) (err error) {
	if (maxDepth != -1 && depth > maxDepth) || !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return visitStructRecursive(v.Elem(), visitor, maxDepth, depth)

	case reflect.Struct:
		err = visitor.BeginStruct(depth, v)
		if err != nil {
			return err
		}
		depth1 := depth + 1
		if maxDepth == -1 || depth1 <= maxDepth {
			t := v.Type()
			n := t.NumField()
			index := 0
			for i := 0; i < n; i++ {
				f := t.Field(i)
				if f.PkgPath == "" && f.Tag.Get("gostart") != "-" { // Only exported fields
					if vi, ok := DereferenceValue(v.Field(i)); ok {
						if f.Anonymous {
							err = visitAnonymousStructFieldRecursive(visitor, vi, &index, depth1)
							if err != nil {
								return err
							}
						} else {
							err = visitor.StructField(depth1, vi, f, index)
							if err != nil {
								return err
							}
							err = visitStructRecursive(vi, visitor, maxDepth, depth1)
							if err != nil {
								return err
							}
							index++
						}
					}
				}
			}
		}
		return visitor.EndStruct(depth, v)

	case reflect.Slice:
		if v.CanSet() {
			modifiedV, err := visitor.PreModifySlice(depth, v)
			if err != nil {
				return err
			}
			v.Set(modifiedV)
		}
		err = visitor.BeginSlice(depth, v)
		if err != nil {
			return err
		}
		depth1 := depth + 1
		if maxDepth == -1 || depth1 <= maxDepth {
			for i := 0; i < v.Len(); i++ {
				if vi, ok := DereferenceValue(v.Index(i)); ok {
					err = visitor.SliceField(depth1, vi, i)
					if err != nil {
						return err
					}
					err = visitStructRecursive(vi, visitor, maxDepth, depth1)
					if err != nil {
						return err
					}
				}
			}
		}
		err = visitor.EndSlice(depth, v)
		if err != nil {
			return err
		}
		if v.CanSet() {
			modifiedV, err := visitor.PostModifySlice(depth, v)
			if err != nil {
				return err
			}
			v.Set(modifiedV)
		}
		return nil

	case reflect.Array:
		err = visitor.BeginArray(depth, v)
		if err != nil {
			return err
		}
		depth1 := depth + 1
		if maxDepth == -1 || depth1 <= maxDepth {
			for i := 0; i < v.Len(); i++ {
				if vi, ok := DereferenceValue(v.Index(i)); ok {
					err = visitor.ArrayField(depth1, vi, i)
					if err != nil {
						return err
					}
					err = visitStructRecursive(vi, visitor, maxDepth, depth1)
					if err != nil {
						return err
					}
				}
			}
		}
		return visitor.EndArray(depth, v)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// ModifySliceStructVisitor

// ModifySliceStructVisitor is a StructVisitor that calls its self function
// value in BeginSlice() and ignores all other StructVisitor methos.
// It can be used to modify the length of slices in complex structs.
type ModifySliceStructVisitor func(depth int, v reflect.Value) (reflect.Value, error)

func (self ModifySliceStructVisitor) BeginStruct(depth int, v reflect.Value) error {
	return nil
}

func (self ModifySliceStructVisitor) StructField(depth int, v reflect.Value, f reflect.StructField, index int) error {
	return nil
}

func (self ModifySliceStructVisitor) EndStruct(depth int, v reflect.Value) error {
	return nil
}

func (self ModifySliceStructVisitor) ModifySlice(depth int, v reflect.Value) (reflect.Value, error) {
	return self(depth, v)
}

func (self ModifySliceStructVisitor) BeginSlice(depth int, v reflect.Value) error {
	return nil
}

func (self ModifySliceStructVisitor) SliceField(depth int, v reflect.Value, index int) error {
	return nil
}

func (self ModifySliceStructVisitor) EndSlice(depth int, v reflect.Value) error {
	return nil
}

func (self ModifySliceStructVisitor) BeginArray(depth int, v reflect.Value) error {
	return nil
}

func (self ModifySliceStructVisitor) ArrayField(depth int, v reflect.Value, index int) error {
	return nil
}

func (self ModifySliceStructVisitor) EndArray(depth int, v reflect.Value) error {
	return nil
}
