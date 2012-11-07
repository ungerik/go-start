package reflection

import (
	"testing"
)

func Test_AssignToResultPtr(t *testing.T) {
	{
		src := 1
		dst := 0
		AssignToResultPtr(src, &dst)
		if dst != 1 {
			t.Fail()
		}
	}

	{
		src := 1
		var dst *int
		AssignToResultPtr(src, &dst)
		if dst == nil || *dst != 1 {
			t.Fail()
		}
	}

	{
		src := new(int)
		*src = 1
		var dst *int
		AssignToResultPtr(src, &dst)
		if dst == nil || *dst != 1 {
			t.Fail()
		}
	}

	{
		src := new(int)
		*src = 1
		dst := 0
		AssignToResultPtr(src, &dst)
		if dst != 1 {
			t.Fail()
		}
	}
}
