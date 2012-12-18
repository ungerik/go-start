package reflection

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
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: X int = 1)
	//   StructField(1: Y float32 = 0)
	//   StructField(2: Z string = "")
	// EndStruct(reflection.exampleStruct)	
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: X int = 0)
	//   StructField(1: Y float32 = 2)
	//   StructField(2: Z string = "")
	// EndStruct(reflection.exampleStruct)	
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
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: X int = 1)
	//   StructField(1: Y float32 = 2)
	//   StructField(2: I int = 0)
	// EndStruct(reflection.exampleStruct)	
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
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: X int = 9)
	//   StructField(1: Slice []int)
	//   BeginSlice([]int)
	//     SliceField(0: int = 1)
	//     SliceField(1: int = 2)
	//   EndSlice([]int)
	//   StructField(2: Array [3]int)
	//   BeginArray([3]int)
	//     ArrayField(0: int = 1)
	//     ArrayField(1: int = 2)
	//     ArrayField(2: int = 3)
	//   EndArray([3]int)
	// EndStruct(reflection.exampleStruct)
}

// func ExampleVisitStruct_flatStructWithMap() {
// 	type exampleStruct struct {
// 		X         int
// 		IntMap    map[string]int
// 		StringMap map[string]string
// 	}
// 	val := &exampleStruct{
// 		X:         9,
// 		IntMap:    map[string]int{"1": 1, "2": 2}, // Reihenfolge in Map nicht definiert
// 		StringMap: map[string]string{"1": "Hello", "2": "World"},
// 	}
// 	VisitStruct(val, NewStdLogStructVisitor())
// 	// Output:
// 	// BeginStruct(reflection.exampleStruct)
// 	//   StructField(0: X int = 9)
// 	//   StructField(1: IntMap map[string]int = map[string]int{"2":2, "1":1})
// 	//   BeginMap(map[string]int)
// 	//     MapField("1": int = 1)
// 	//     MapField("2": int = 2)
// 	//   EndMap(map[string]int)
// 	//   StructField(2: StringMap map[string]string = map[string]string{"2":"World", "1":"Hello"})
// 	//   BeginMap(map[string]string)
// 	//     MapField("1": string = "Hello")
// 	//     MapField("2": string = "World")
// 	//   EndMap(map[string]string)
// 	// EndStruct(reflection.exampleStruct)
// }

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
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: A0 int = 0)
	//   StructField(1: A1 int = 1)
	//   StructField(2: A2 int = 2)
	//   StructField(3: B0 bool = false)
	//   StructField(4: B1 bool = true)
	// EndStruct(reflection.exampleStruct)
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
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: A0 int = 0)
	//   StructField(1: A1 int = 1)
	//   StructField(2: C0 string = "C0")
	//   StructField(3: B0 bool = false)
	//   StructField(4: A0 int = 0)
	//   StructField(5: A1 int = 1)
	//   StructField(6: B1 bool = true)
	//   StructField(7: C1 string = "C1")
	// EndStruct(reflection.exampleStruct)
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
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: C reflection.C)
	//   BeginStruct(reflection.C)
	//     StructField(0: A reflection.A)
	//     BeginStruct(reflection.A)
	//       StructField(0: A0 int = 0)
	//       StructField(1: A1 int = 1)
	//     EndStruct(reflection.A)
	//     StructField(1: C0 string = "C0")
	//     StructField(2: B reflection.B)
	//     BeginStruct(reflection.B)
	//       StructField(0: B0 bool = false)
	//       StructField(1: A reflection.A)
	//       BeginStruct(reflection.A)
	//         StructField(0: A0 int = 0)
	//         StructField(1: A1 int = 1)
	//       EndStruct(reflection.A)
	//       StructField(2: B1 bool = true)
	//     EndStruct(reflection.B)
	//     StructField(3: C1 string = "C1")
	//   EndStruct(reflection.C)
	// EndStruct(reflection.exampleStruct)
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
	// BeginStruct(reflection.exampleStruct)
	//   StructField(0: C reflection.C)
	//   BeginStruct(reflection.C)
	//     StructField(0: A reflection.A)
	//     BeginStruct(reflection.A)
	//     EndStruct(reflection.A)
	//     StructField(1: C0 string = "C0")
	//     StructField(2: B reflection.B)
	//     BeginStruct(reflection.B)
	//     EndStruct(reflection.B)
	//     StructField(3: C1 string = "C1")
	//   EndStruct(reflection.C)
	// EndStruct(reflection.exampleStruct)
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
	// BeginStruct(reflection.person)
	//   StructField(0: Name string = "Erik Unger")
	//   StructField(1: Email []reflection.emailIdentity)
	//   BeginSlice([]reflection.emailIdentity)
	//     SliceField(0: reflection.emailIdentity)
	//     BeginStruct(reflection.emailIdentity)
	//       StructField(0: Address string = "erik@erikunger.com")
	//       StructField(1: Description string = "Test")
	//     EndStruct(reflection.emailIdentity)
	//   EndSlice([]reflection.emailIdentity)
	//   StructField(2: ExtraInfo string = "info")
	// EndStruct(reflection.person)
}
