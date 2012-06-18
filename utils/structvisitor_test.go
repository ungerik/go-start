package utils

import (
// "fmt"
// "testing"
)

func ExampleVisitStruct_simpleFlatStruct() {
	type exampleStruct struct {
		X      int
		Y      float32
		Z      string
		hidden string
	}
	VisitStruct(&exampleStruct{X: 1}, NewStdLogStructVisitor())
	VisitStruct(exampleStruct{Y: 2}, NewStdLogStructVisitor())
	// Output:
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, X int = 1)
	//   StructField(1, Y float32 = 0)
	//   StructField(2, Z string = "")
	// EndStruct(utils.exampleStruct)	
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, X int = 0)
	//   StructField(1, Y float32 = 2)
	//   StructField(2, Z string = "")
	// EndStruct(utils.exampleStruct)	
}

func ExampleVisitStruct_simpleFlatPtrsStruct() {
	type exampleStruct struct {
		X *int
		Y **float32
		Z *string
		N **string
		I interface{}
	}
	val := &exampleStruct{
		X: new(int),
		Y: new(*float32),
		Z: nil,
		N: new(*string),
		I: new(int),
	}
	*val.X = 1
	*val.Y = new(float32)
	**val.Y = 2
	VisitStruct(val, NewStdLogStructVisitor())
	// Output:
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, X int = 1)
	//   StructField(1, Y float32 = 2)
	//   StructField(2, I int = 0)
	// EndStruct(utils.exampleStruct)	
}

func ExampleVisitStruct_flatStructWithSliceAndArray() {
	type exampleStruct struct {
		X      int
		hidden bool
		Slice  []int
		Array  [3]int
	}
	val := &exampleStruct{
		X:     9,
		Slice: []int{1, 2},
		Array: [...]int{1, 2, 3},
	}
	VisitStruct(val, NewStdLogStructVisitor())
	// Output:
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, X int = 9)
	//   StructField(1, Slice []int)
	//   BeginSlice([]int)
	//     SliceField(0, int = 1)
	//     SliceField(1, int = 2)
	//   EndSlice([]int)
	//   StructField(2, Array [3]int)
	//   BeginArray([3]int)
	//     ArrayField(0, int = 1)
	//     ArrayField(1, int = 2)
	//     ArrayField(2, int = 3)
	//   EndArray([3]int)
	// EndStruct(utils.exampleStruct)
}

func ExampleVisitStruct_flatAnonymousStructFields() {
	type A struct {
		A0 int
		A1 interface{}
		A2 int
	}
	type B struct {
		B0 bool
		B1 bool
	}
	type exampleStruct struct {
		A
		*B
	}
	val := &exampleStruct{
		A: A{A0: 0, A1: 1, A2: 2},
		B: &B{B0: false, B1: true},
	}
	VisitStruct(val, NewStdLogStructVisitor())
	// Output:
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, A0 int = 0)
	//   StructField(1, A1 int = 1)
	//   StructField(2, A2 int = 2)
	//   StructField(3, B0 bool = false)
	//   StructField(4, B1 bool = true)
	// EndStruct(utils.exampleStruct)
}

func ExampleVisitStruct_deepAnonymousStructFields() {
	type A struct {
		A0 interface{}
		A1 int
	}
	type B struct {
		B0 bool
		A
		B1 bool
	}
	type C struct {
		A
		C0 string
		*B
		C1 string
	}
	type exampleStruct struct {
		C
	}
	val := &exampleStruct{
		C: C{
			A:  A{A0: 0, A1: 1},
			C0: "C0",
			B: &B{
				B0: false,
				A:  A{A0: 0, A1: 1},
				B1: true,
			},
			C1: "C1",
		},
	}
	VisitStruct(val, NewStdLogStructVisitor())
	// Output:
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, A0 int = 0)
	//   StructField(1, A1 int = 1)
	//   StructField(2, C0 string = "C0")
	//   StructField(3, B0 bool = false)
	//   StructField(4, A0 int = 0)
	//   StructField(5, A1 int = 1)
	//   StructField(6, B1 bool = true)
	//   StructField(7, C1 string = "C1")
	// EndStruct(utils.exampleStruct)
}

func ExampleVisitStruct_deepStructFields() {
	type A struct {
		A0 interface{}
		A1 int
	}
	type B struct {
		B0 bool
		A  A
		B1 bool
	}
	type C struct {
		A  A
		C0 string
		B  *B
		C1 string
	}
	type exampleStruct struct {
		C C
	}
	val := &exampleStruct{
		C: C{
			A:  A{A0: 0, A1: 1},
			C0: "C0",
			B: &B{
				B0: false,
				A:  A{A0: 0, A1: 1},
				B1: true,
			},
			C1: "C1",
		},
	}
	VisitStruct(val, NewStdLogStructVisitor())
	// Output:
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, C utils.C)
	//   BeginStruct(utils.C)
	//     StructField(0, A utils.A)
	//     BeginStruct(utils.A)
	//       StructField(0, A0 int = 0)
	//       StructField(1, A1 int = 1)
	//     EndStruct(utils.A)
	//     StructField(1, C0 string = "C0")
	//     StructField(2, B utils.B)
	//     BeginStruct(utils.B)
	//       StructField(0, B0 bool = false)
	//       StructField(1, A utils.A)
	//       BeginStruct(utils.A)
	//         StructField(0, A0 int = 0)
	//         StructField(1, A1 int = 1)
	//       EndStruct(utils.A)
	//       StructField(2, B1 bool = true)
	//     EndStruct(utils.B)
	//     StructField(3, C1 string = "C1")
	//   EndStruct(utils.C)
	// EndStruct(utils.exampleStruct)
}

func ExampleVisitStruct_limitDepth() {
	type A struct {
		A0 interface{}
		A1 int
	}
	type B struct {
		B0 bool
		A  A
		B1 bool
	}
	type C struct {
		A  A
		C0 string
		B  *B
		C1 string
	}
	type exampleStruct struct {
		C C
	}
	val := &exampleStruct{
		C: C{
			A:  A{A0: 0, A1: 1},
			C0: "C0",
			B: &B{
				B0: false,
				A:  A{A0: 0, A1: 1},
				B1: true,
			},
			C1: "C1",
		},
	}
	VisitStructDepth(val, NewStdLogStructVisitor(), 2)
	// Output:
	// BeginStruct(utils.exampleStruct)
	//   StructField(0, C utils.C)
	//   BeginStruct(utils.C)
	//     StructField(0, A utils.A)
	//     BeginStruct(utils.A)
	//     EndStruct(utils.A)
	//     StructField(1, C0 string = "C0")
	//     StructField(2, B utils.B)
	//     BeginStruct(utils.B)
	//     EndStruct(utils.B)
	//     StructField(3, C1 string = "C1")
	//   EndStruct(utils.C)
	// EndStruct(utils.exampleStruct)
}

func ExampleVisitStruct_sliceOfStructInAnonymousStruct() {
	type emailIdentity struct {
		Address     string
		Description string
	}
	type user struct {
		Name  string
		Email []emailIdentity
	}
	type person struct {
		user
		ExtraInfo string
	}
	val := &person{
		user: user{
			Name: "Erik Unger",
			Email: []emailIdentity{
				emailIdentity{Address: "erik@erikunger.com", Description: "Test"},
			},
		},
		ExtraInfo: "info",
	}
	VisitStruct(val, NewStdLogStructVisitor())
	// Output:
	// BeginStruct(utils.person)
	//   StructField(0, Name string = "Erik Unger")
	//   StructField(1, Email []utils.emailIdentity)
	//   BeginSlice([]utils.emailIdentity)
	//     SliceField(0, utils.emailIdentity)
	//     BeginStruct(utils.emailIdentity)
	//       StructField(0, Address string = "erik@erikunger.com")
	//       StructField(1, Description string = "Test")
	//     EndStruct(utils.emailIdentity)
	//   EndSlice([]utils.emailIdentity)
	//   StructField(2, ExtraInfo string = "info")
	// EndStruct(utils.person)
}
