package templatesystem

import (
	"fmt"
	"github.com/ungerik/go-start/errs"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

type PrintfTemplate struct {
	text string
}

func (self *PrintfTemplate) Render(out io.Writer, context interface{}) (err error) {
	var str string
	switch reflect.TypeOf(context).Kind() {
	case reflect.Slice, reflect.Array:
		panic("todo implementation")
	default:
		str = fmt.Sprintf(self.text, context)
	}
	if strings.HasPrefix(str, "%!") {
		return errs.Format(str)
	}
	_, err = out.Write([]byte(str))
	return err
}

type Printf struct{}

func (self *Printf) ParseFile(filename string) (Template, error) {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return self.ParseString(string(text), filepath.Base(filename))
}

func (self *Printf) ParseString(text, name string) (Template, error) {
	return &PrintfTemplate{text: text}, nil
}
