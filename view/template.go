package view

import (
	"path"
	"time"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/templatesystem"
)

type GetTemplateContextFunc func(ctx *Context) (context interface{}, err error)

func NewTemplate(filename string, getContext GetTemplateContextFunc) *Template {
	return &Template{Filename: filename, GetContext: getContext}
}

func TemplateContext(context interface{}) GetTemplateContextFunc {
	return func(ctx *Context) (interface{}, error) {
		return context, nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// Template

type Template struct {
	ViewWithURLBase
	Filename       string // Will set file extension at ContentType
	Text           string
	ContentTypeExt string
	GetContext     GetTemplateContextFunc
	TemplateSystem templatesystem.Implementation // If nil, self.App.Config.TemplateSystem is used
	template       templatesystem.Template
	modifiedTime   time.Time
}

func (self *Template) parseTemplate() (templ templatesystem.Template, err error) {
	templateSystem := self.TemplateSystem
	if templateSystem == nil {
		templateSystem = Config.TemplateSystem
	}

	if self.Filename == "" {
		self.modifiedTime = time.Time{}
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

func (self *Template) Render(ctx *Context) (err error) {
	if self.template != nil && self.Filename != "" {
		_, found, modified := FindTemplateFile(self.Filename)
		if !found {
			return errs.Format("Template file not found: %s", self.Filename)
		}
		if modified.After(self.modifiedTime) {
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
		ctx.Response.SetContentTypeByExt(self.ContentTypeExt)
	}

	var context interface{}
	if self.GetContext != nil {
		context, err = self.GetContext(ctx)
		if err != nil {
			return err
		}
	}

	// todo: how to add config data to context if it's not a slice?
	// map[string][]string{"args": ctx.Request.URLArgs}
	// Config, context.Web
	return self.template.Render(ctx.Response, context)
}
