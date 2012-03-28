package templatesystem

import (
	"github.com/ungerik/mustache.go"
	"io"
)

type MustacheTemplate struct {
	templ *mustache.Template
}

func (self *MustacheTemplate) Render(out io.Writer, context interface{}) (err error) {
	str := self.templ.Render(context)
	_, err = out.Write([]byte(str))
	return err
}

type Mustache struct{}

func (self *Mustache) ParseFile(filename string) (Template, error) {
	templ, err := mustache.ParseFile(filename)
	if err != nil {
		return nil, err
	}
	return &MustacheTemplate{templ: templ}, nil
}

func (self *Mustache) ParseString(text, name string) (Template, error) {
	templ, err := mustache.ParseString(text)
	if err != nil {
		return nil, err
	}
	return &MustacheTemplate{templ: templ}, nil
}
