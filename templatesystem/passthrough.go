package templatesystem

import (
	"io"
	"io/ioutil"
	"path/filepath"
)

type PassThroughTemplate struct {
	text string
}

func (self *PassThroughTemplate) Render(out io.Writer, context interface{}) (err error) {
	_, err = out.Write([]byte(self.text))
	return err
}

type PassThrough struct{}

func (self *PassThrough) ParseFile(filename string) (Template, error) {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return self.ParseString(string(text), filepath.Base(filename))
}

func (self *PassThrough) ParseString(text, name string) (Template, error) {
	return &PassThroughTemplate{text: text}, nil
}
