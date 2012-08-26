package reflection

import "reflect"

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

func (self ModifySliceStructVisitor) BeginMap(depth int, v reflect.Value) error {
	return nil
}

func (self ModifySliceStructVisitor) MapField(depth int, v reflect.Value, key string, index int) error {
	return nil
}

func (self ModifySliceStructVisitor) EndMap(depth int, v reflect.Value) error {
	return nil
}
