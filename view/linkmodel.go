package view

import (
	"fmt"
	"reflect"
)

///////////////////////////////////////////////////////////////////////////////
// LinkModel

type LinkModel interface {
	URL
	LinkContent(context *Context) View
	LinkTitle(context *Context) string
	LinkRel(context *Context) string
}

func NewLinkModel(url interface{}, content ...interface{}) LinkModel {
	getContent := func() View {
		if len(content) == 0 {
			return NewView(url)
		}
		return NewViews(content...)
	}
	switch s := url.(type) {
	case **Page:
		return &PageLink{Page: s, Content: getContent()}
	case *ViewWithURL:
		return &URLLink{Url: &indirectViewWithURL{s}, Content: NewViews(content...)}
	case LinkModel:
		if len(content) > 0 {
			return &URLLink{Url: s, Content: NewViews(content...)}
		}
		return s
	case URL:
		return &URLLink{Url: s, Content: getContent()}
	case fmt.Stringer:
		return &StringLink{Url: s.String(), Content: getContent()}
	}
	v := reflect.ValueOf(url)
	if v.Kind() != reflect.String {
		panic(fmt.Errorf("Invalid type for url: %T", url))
	}
	return &StringLink{Url: v.String(), Content: getContent()}
}

func NewLinkModelRel(url interface{}, rel string, content ...interface{}) LinkModel {
	getContent := func() View {
		if len(content) == 0 {
			return NewView(url)
		}
		return NewViews(content...)
	}
	switch s := url.(type) {
	case **Page:
		return &PageLink{Page: s, Content: getContent(), Rel: rel}
	case *ViewWithURL:
		return &URLLink{Url: &indirectViewWithURL{s}, Content: NewViews(content...), Rel: rel}
	case LinkModel:
		if len(content) > 0 {
			return &URLLink{Url: s, Content: NewViews(content...), Rel: rel}
		}
		return s
	case URL:
		return &URLLink{Url: s, Content: getContent(), Rel: rel}
	case fmt.Stringer:
		return &StringLink{Url: s.String(), Content: getContent(), Rel: rel}
	}
	v := reflect.ValueOf(url)
	if v.Kind() != reflect.String {
		panic(fmt.Errorf("Invalid type for url: %T", url))
	}
	return &StringLink{Url: v.String(), Content: getContent(), Rel: rel}
}

///////////////////////////////////////////////////////////////////////////////
// PageLink

func NewPageLink(page **Page, title string) *PageLink {
	return &PageLink{Page: page, Title: title}
}

type PageLink struct {
	Page    **Page
	Content View   // If is nil, then self.LinkTitle() will be used
	Title   string // If is "", then self.Page.LinkTitle(context) will be used
	Rel     string
}

func (self *PageLink) URL(context *Context, args ...string) string {
	return self.Page.URL(context, args...)
}

func (self *PageLink) LinkContent(context *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(context))
	}
	return self.Content
}

func (self *PageLink) LinkTitle(context *Context) string {
	if self.Title == "" {
		return self.Page.LinkTitle(context)
	}
	return self.Title
}

func (self *PageLink) LinkRel(context *Context) string {
	return self.Rel
}

///////////////////////////////////////////////////////////////////////////////
// StringLink

type StringLink struct {
	Url     string
	Content View   // If is nil, then self.LinkTitle() will be used
	Title   string // If is "", then self.URL will be used
	Rel     string
}

func (self *StringLink) URL(context *Context, args ...string) string {
	return StringURL(self.Url).URL(context, args...)
}

func (self *StringLink) LinkContent(context *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(context))
	}
	return self.Content
}

func (self *StringLink) LinkTitle(context *Context) string {
	if self.Title == "" {
		return self.Url
	}
	return self.Title
}

func (self *StringLink) LinkRel(context *Context) string {
	return self.Rel
}

///////////////////////////////////////////////////////////////////////////////
// URLLink

type URLLink struct {
	Url     URL
	Content View   // If is nil, then self.LinkTitle() will be used
	Title   string // If is "", then self.URL will be used
	Rel     string
}

func (self *URLLink) URL(context *Context, args ...string) string {
	return self.Url.URL(context, args...)
}

func (self *URLLink) LinkContent(context *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(context))
	}
	return self.Content
}

func (self *URLLink) LinkTitle(context *Context) string {
	if self.Title == "" {
		return self.Url.URL(context)
	}
	return self.Title
}

func (self *URLLink) LinkRel(context *Context) string {
	return self.Rel
}
