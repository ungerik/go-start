package templatesystem

/*
import (
	"io"
	"io/ioutil"
	"kasia"
	"path/filepath"
)

type KasiaTemplate struct {
	templ *kasia.Template
}

func (self *KasiaTemplate) Render(out io.Writer, context interface{}) (err error) {
	return self.templ.Run(out, context)
}

type Kasia struct{}

func (self *Kasia) ParseFile(filename string) (Template, error) {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return self.ParseString(string(text), filepath.Base(filename))
}

func (self *Kasia) ParseString(text, name string) (Template, error) {
	templ := kasia.New()
	err := templ.Parse([]byte(text))
	if err != nil {
		return nil, err
	}
	return &KasiaTemplate{templ: templ}, nil
}
*/
