package reflection

import (
	"fmt"
	"reflect"
	"sort"
)

// Sort will sort slice according to compareFunc using reflection.
// slice can be a slice of any element type including interface{}.
// compareFunc must have two arguments that are assignable from
// the slice element type or pointers to such a type.
// The result of compareFunc must be a bool indicating
// if the first argument is less than the second.
// If the element type of slice is interface{}, then the type
// of the compareFunc arguments can be any type and dynamic
// casting from the interface value or its address will be attempted.
func Sort(slice, compareFunc interface{}) {
	f, err := NewSortCompareFunc(compareFunc)
	if err != nil {
		panic(err)
	}
	err = f.Sort(slice)
	if err != nil {
		panic(err)
	}
}

func NewSortCompareFunc(compareFunc interface{}) (*SortCompareFunc, error) {
	t := reflect.TypeOf(compareFunc)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("compareFunc must be a function, got %T", compareFunc)
	}
	if t.NumIn() != 2 {
		return nil, fmt.Errorf("compareFunc must take two arguments, got %d", t.NumIn())
	}
	if t.In(0) != t.In(1) {
		return nil, fmt.Errorf("compareFunc's arguments must be identical, got %s and %s", t.In(0), t.In(1))
	}
	if t.NumOut() != 1 {
		return nil, fmt.Errorf("compareFunc must have one result, got %d", t.NumOut())
	}
	if t.Out(0).Kind() != reflect.Bool {
		return nil, fmt.Errorf("compareFunc result must be bool, got %s", t.Out(0))
	}

	argType := t.In(0)
	argTypePtr := argType.Kind() == reflect.Ptr
	if argTypePtr {
		argType = argType.Elem()
	}

	return &SortCompareFunc{
		CompareFunc: reflect.ValueOf(compareFunc),
		ArgType:     argType,
		PtrArgs:     argTypePtr,
	}, nil
}

type SortCompareFunc struct {
	CompareFunc reflect.Value
	ArgType     reflect.Type
	PtrArgs     bool
}

func (self *SortCompareFunc) Sort(slice interface{}) error {
	sliceV := reflect.ValueOf(slice)
	if sliceV.Kind() != reflect.Slice {
		return fmt.Errorf("Need slice got %T", slice)
	}
	elemT := sliceV.Type().Elem()
	if elemT != self.ArgType && elemT.Kind() != reflect.Interface {
		return fmt.Errorf("Slice element type must be interface{} or %s, got %s", self.ArgType, elemT)
	}
	sort.Sort(&sortable{self, sliceV})
	return nil
}

// Implements sort.Interface
type sortable struct {
	*SortCompareFunc
	slice reflect.Value
}

func (self *sortable) Len() int {
	return self.slice.Len()
}

func (self *sortable) Less(i, j int) bool {
	arg0 := self.slice.Index(i)
	arg1 := self.slice.Index(j)
	if self.slice.Type().Elem().Kind() == reflect.Interface {
		arg0 = arg0.Elem()
		arg1 = arg1.Elem()
	}
	if (arg0.Kind() == reflect.Ptr) != self.PtrArgs {
		if self.PtrArgs {
			// Expects PtrArgs for SortCompareFunc, but slice is value type
			arg0 = arg0.Addr()
		} else {
			// Expects value type for SortCompareFunc, but slice is PtrArgs
			arg0 = arg0.Elem()
		}
	}
	if (arg1.Kind() == reflect.Ptr) != self.PtrArgs {
		if self.PtrArgs {
			// Expects PtrArgs for SortCompareFunc, but slice is value type
			arg1 = arg1.Addr()
		} else {
			// Expects value type for SortCompareFunc, but slice is PtrArgs
			arg1 = arg1.Elem()
		}
	}
	return self.CompareFunc.Call([]reflect.Value{arg0, arg1})[0].Bool()
}

func (self *sortable) Swap(i, j int) {
	temp := self.slice.Index(i).Interface()
	self.slice.Index(i).Set(self.slice.Index(j))
	self.slice.Index(j).Set(reflect.ValueOf(temp))
}
