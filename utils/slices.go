package utils

import (
	// "github.com/ungerik/go-start/errs"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func StringSliceUnion(a, b []string) []string {
	if len(b) == 0 {
		return a
	}
	if len(a) == 0 {
		return b
	}
	result := a
	for _, sb := range b {
		found := false
		for _, sa := range a {
			if sa == sb {
				found = true
				break
			}
		}
		if !found {
			result = append(result, sb)
		}
	}
	return result
}

func AppendEmptySliceField(slice reflect.Value) reflect.Value {
	newField := reflect.Zero(slice.Type().Elem())
	return reflect.Append(slice, newField)
}

// Sets the length of a slice by sub-slicing a slice that's too long,
// or appending empty fields if slice is too short.
func SetSliceLengh(slice reflect.Value, length int) reflect.Value {
	if length > slice.Len() {
		for i := slice.Len(); i < length; i++ {
			slice = AppendEmptySliceField(slice)
		}
	} else if length < slice.Len() {
		slice = slice.Slice(0, length)
	}
	return slice
}

func DeleteEmptySliceElementsVal(sliceVal reflect.Value) reflect.Value {
	if sliceVal.Kind() != reflect.Slice {
		panic("Argument is not a slice: " + sliceVal.String())
	}
	zeroVal := reflect.Zero(sliceVal.Type().Elem())
	for i := 0; i < sliceVal.Len(); i++ {
		elemVal := sliceVal.Index(i)
		if reflect.DeepEqual(elemVal.Interface(), zeroVal.Interface()) {
			before := sliceVal.Slice(0, i)
			after := sliceVal.Slice(i+1, sliceVal.Len())
			sliceVal = reflect.AppendSlice(before, after)
			i--
		}
	}
	return sliceVal
}

func DeleteEmptySliceElements(slice interface{}) interface{} {
	return DeleteEmptySliceElementsVal(reflect.ValueOf(slice)).Interface()
}

// func SliceInsert(slice []interface{}, index int, count int, value interface{}) (result []interface{}) {
// 	switch {
// 	case count < 0:
// 		panic(errs.Format("Negative count %d not allowed", count))
// 	case count == 0:
// 		return slice
// 	}

// 	length := len(slice)
// 	errs.PanicIfErrIndexOutOfBounds("SliceInsert", index, length)

// 	result = make([]interface{}, length+count)
// 	copy(result, slice[:index])
// 	copy(result[index+count:], slice[index:])
// 	for i := index; i < index+count; i++ {
// 		result[i] = value
// 	}

// 	return result
// }

// func SliceDelete(slice []interface{}, index int, count int) (result []interface{}) {
// 	switch {
// 	case count < 0:
// 		panic(errs.Format("Negative count %d not allowed", count))
// 	case count == 0:
// 		return slice
// 	}

// 	length := len(slice)
// 	errs.PanicIfErrIndexOutOfBounds("SliceDelete", index, length)

// 	if index+count > length {
// 		count = length - index
// 	}

// 	return append(slice[:index], slice[index+count:]...)
// }

// Implements sort.Interface
type SortableInterfaceSlice struct {
	Slice    []interface{}
	LessFunc func(a, b interface{}) bool
}

func (self *SortableInterfaceSlice) Len() int {
	return len(self.Slice)
}

func (self *SortableInterfaceSlice) Less(i, j int) bool {
	return self.LessFunc(self.Slice[i], self.Slice[j])
}

func (self *SortableInterfaceSlice) Swap(i, j int) {
	self.Slice[i], self.Slice[j] = self.Slice[j], self.Slice[i]
}

func (self *SortableInterfaceSlice) Sort() {
	sort.Sort(self)
}

func SortInterfaceSlice(slice []interface{}, lessFunc func(a, b interface{}) bool) {
	sortable := SortableInterfaceSlice{slice, lessFunc}
	sortable.Sort()
}

/*
func CloneStringSlice(original []string) (clone []string) {
	if original != nil {
		clone = make([]string, len(original))
		for i := range original {
			clone[i] = original[i]
		}
	}
	return clone
}


func CloneByteSlice(original []byte) (clone []byte) {
	if original != nil {
		clone = make([]byte, len(original))
		for i := range original {
			clone[i] = original[i]
		}
	}
	return clone
}
*/

func MakeVersionTuple(fields ...int) VersionTuple {
	t := make(VersionTuple, len(fields))
	for i := range fields {
		t[i] = fields[i]
	}
	return t
}

func ParseVersionTuple(s string) (VersionTuple, error) {
	fields := strings.Split(s, ".")
	t := make(VersionTuple, len(fields))
	for i := range fields {
		value, err := strconv.ParseInt(fields[i], 10, 32)
		if err != nil {
			return nil, err
		}
		t[i] = int(value)
	}
	return t, nil
}

type VersionTuple []int

func (self VersionTuple) GreaterEqual(other VersionTuple) bool {
	for i := range other {
		var value int
		if i < len(self) {
			value = self[i]
		}
		if value > other[i] {
			return true
		} else if value < other[i] {
			return false
		}
	}
	return true
}

func (self VersionTuple) String() string {
	var sb StringBuilder
	for i := range self {
		if i > 0 {
			sb.Byte('.')
		}
		sb.Int(self[i])
	}
	return sb.String()
}
