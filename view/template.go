package view

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/templatesystem"
	"github.com/ungerik/go-start/utils"
	"path"
)

type GetTemplateContextFunc func(requestContext *Context) (context interface{}, err error)

func NewTemplate(filename string, getContext GetTemplateContextFunc) *Template {
	return &Template{Filename: filename, GetContext: getContext}
}

func NewHTML5BoilerplateCSSTemplate(getContext GetTemplateContextFunc, filenames ...string) Views {
	views := make(Views, len(filenames)+2)
	views[0] = &StaticFile{Filename: "css/html5boilerplate/normalize.css"}
	for i := range filenames {
		views[i+1] = NewTemplate(filenames[i], getContext)
	}
	views[len(views)-1] = &StaticFile{Filename: "css/html5boilerplate/poststyle.css"}
	return views
}

func TemplateContext(context interface{}) GetTemplateContextFunc {
	return func(requestContext *Context) (interface{}, error) {
		return context, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// Template

type Template struct {
	ViewBase
	Filename       string // Will set file extension at ContentType
	Text           string
	ContentTypeExt string
	GetContext     GetTemplateContextFunc
	TemplateSystem templatesystem.Implementation // If nil, self.App.Config.TemplateSystem is used
	template       templatesystem.Template
	modifiedTime   int64
}

func (self *Template) parseTemplate() (templ templatesystem.Template, err error) {
	templateSystem := self.TemplateSystem
	if templateSystem == nil {
		templateSystem = Config.TemplateSystem
	}

	if self.Filename == "" {
		self.modifiedTime = 0
		return templateSystem.ParseString(self.Text, "")
	}

	filePath, found, modified := FindTemplateFile(self.Filename)
	if !found {
		return nil, errs.Format("Template file not found: %s", self.Filename)
	}
	self.modifiedTime = modified

	if self.ContentTypeExt == "" {
		self.ContentTypeExt = path.Ext(filePath)
	}

	return templateSystem.ParseFile(filePath)
}

func (self *Template) Render(requestContext *Context, writer *utils.XMLWriter) (err error) {
	if self.template != nil && self.Filename != "" {
		_, found, modified := FindTemplateFile(self.Filename)
		if !found {
			return errs.Format("Template file not found: %s", self.Filename)
		}
		if modified > self.modifiedTime {
			self.template = nil
			self.ContentTypeExt = ""
		}
	}

	if self.template == nil {
		self.template, err = self.parseTemplate()
		if err != nil {
			return err
		}
	}

	if self.ContentTypeExt != "" {
		requestContext.ContentType(self.ContentTypeExt)
	}

	var context interface{}
	if self.GetContext != nil {
		context, err = self.GetContext(requestContext)
		if err != nil {
			return err
		}
	}

	// todo: how to add config data to context if it's not a slice?
	// map[string][]string{"args": context.PathArgs}
	// Config, context.Web
	return self.template.Render(writer, context)
}

//func (self *Template) SetFilename(filename string) {
//	self.Filename = filename
//	self.Text = ""
//	self.template = nil
//	ViewChanged(self)
//}
//
//func (self *Template) SetText(text string) {
//	self.Filename = ""
//	self.Text = text
//	self.template = nil
//	ViewChanged(self)
//}
//
//func (self *Template) SetContext(context interface{}) {
//	self.Context = context
//	ViewChanged(self)
//}
