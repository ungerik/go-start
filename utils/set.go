package utils

type dummy struct {}

type StringSet struct {
	m map[string]dummy
}

func (self StringSet) Add(val string) {
	self.m[val] = dummy{}
}

func (self StringSet) Delete(val string) {
	delete(self.m, val)
}

func (self StringSet) Contains(val string) bool {
	_, ok := self.m[val]
	return ok
}

type IntSet struct {
	m map[int]dummy
}

func (self IntSet) Add(val int) {
	self.m[val] = dummy{}
}

func (self IntSet) Delete(val int) {
	delete(self.m, val)
}

func (self IntSet) Contains(val int) bool {
	_, ok := self.m[val]
	return ok
}
