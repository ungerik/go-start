package view

type Redirect string

func (self Redirect) Error() string {
	return string(self)
}

type PermanentRedirect string

func (self PermanentRedirect) Error() string {
	return string(self)
}

type NotFound string

func (self NotFound) Error() string {
	return string(self)
}

type Forbidden string

func (self Forbidden) Error() string {
	return string(self)
}
