package reflection

import (
	"reflect"

	// "github.com/ungerik/go-start/debug"
)

////////////////////////////////////////////////////////////////////////////////
// StructVisitor

type StructVisitor interface {
	BeginStruct(depth int, v reflect.Value) error
	StructField(depth int, v reflect.Value, f reflect.StructField, index int) error
	EndStruct(depth int, v reflect.Value) error

	BeginSlice(depth int, v reflect.Value) error
	SliceField(depth int, v reflect.Value, index int) error
	EndSlice(depth int, v reflect.Value) error

	BeginArray(depth int, v reflect.Value) error
	ArrayField(depth int, v reflect.Value, index int) error
	EndArray(depth int, v reflect.Value) error

	BeginMap(depth int, v reflect.Value) error
	MapField(depth int, v reflect.Value, key string, index int) error
	EndMap(depth int, v reflect.Value) error
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

func visitAnonymousStructFieldRecursive(visitor StructVisitor, v reflect.Value, maxDepth, depth int, index *int) (err error) {
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath == "" && f.Tag.Get("gostart") != "-" { // Only exported fields
				if vi, ok := DereferenceValue(v.Field(i)); ok {
					if f.Anonymous {
						err = visitAnonymousStructFieldRecursive(visitor, vi, maxDepth, depth, index)
						if err != nil {
							return err
						}
					} else {
						err = visitor.StructField(depth, vi, f, *index)
						if err != nil {
							return err
						}
						err = visitStructRecursive(vi, visitor, maxDepth, depth)
						if err != nil {
							return err
						}
						*index++
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
			index := 0
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.PkgPath == "" && f.Tag.Get("gostart") != "-" { // Only exported fields
					if vi, ok := DereferenceValue(v.Field(i)); ok {
						if f.Anonymous {
							err = visitAnonymousStructFieldRecursive(visitor, vi, maxDepth, depth1, &index)
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

	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String && v.Len() > 0 {
			err = visitor.BeginMap(depth, v)
			if err != nil {
				return err
			}
			depth1 := depth + 1
			if maxDepth == -1 || depth1 <= maxDepth {
				for i, key := range v.MapKeys() {
					if vi, ok := DereferenceValue(v.MapIndex(key)); ok {
						err = visitor.MapField(depth1, vi, key.String(), i)
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
			return visitor.EndMap(depth, v)
		}

	case reflect.Slice:
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
		return visitor.EndSlice(depth, v)

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
