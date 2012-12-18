package reflection

import (
	"testing"
)

func Test_SmartCopy(t *testing.T) {
	{
		src := 1
		dst := 0
		SmartCopy(src, &dst)
		if dst != 1 {
			t.Fail()
		}
	}
	{
		src := 1
		var dst *int
		SmartCopy(src, &dst)
		if dst == nil || *dst != 1 {
			t.Fail()
		}
	}
	return // todo

	{
		src := new(int)
		*src = 1
		var dst *int
		SmartCopy(src, &dst)
		if dst == nil || *dst != 1 {
			t.Fail()
		}
	}

	{
		src := new(int)
		*src = 1
		dst := 0
		SmartCopy(src, &dst)
		if dst != 1 {
			t.Fail()
		}
	}
}
