package templatesystem

import (
	"io"
	"text/template"
)

type GoTemplate struct {
	templ *template.Template
}

func (self *GoTemplate) Render(out io.Writer, context interface{}) (err error) {
	return self.templ.Execute(out, context)
}

type Go struct{}

func (self *Go) ParseFile(filename string) (Template, error) {
	templ, err := template.ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	return &GoTemplate{templ}, nil
}

func (self *Go) ParseString(text, name string) (Template, error) {
	templ, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}
	return &GoTemplate{templ}, nil
}
