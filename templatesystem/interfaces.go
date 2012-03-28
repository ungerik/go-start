package templatesystem

import "io"

type Template interface {
	Render(out io.Writer, context interface{}) error
}

type Implementation interface {
	ParseFile(filename string) (Template, error)
	ParseString(text, name string) (Template, error)
}
