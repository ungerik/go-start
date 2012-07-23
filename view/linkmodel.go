package view

import (
	"fmt"
	"reflect"
)

type LinkModel interface {
	URL
	LinkContent(response *Response) View
	LinkTitle(response *Response) string
	LinkRel(response *Response) string
}

func NewLinkModel(url interface{}, content ...interface{}) LinkModel {
	getContent := func() View {
		if len(content) == 0 {
			return NewView(url)
		}
		return NewViews(content...)
	}
	getContentString := func() string {
		if len(content) == 0 {
			return ""
		}
		s, _ := content[0].(string)
		return s
	}
	switch s := url.(type) {
	case **Page:
		return &PageLink{Page: s, Content: getContent(), Title: getContentString()}
	case *ViewWithURL:
		return &URLLink{Url: IndirectViewWithURL(s), Content: NewViews(content...), Title: getContentString()}
	case LinkModel:
		if len(content) > 0 {
			return &URLLink{Url: s, Content: NewViews(content...), Title: getContentString()}
		}
		return s
	case URL:
		return &URLLink{Url: s, Content: getContent(), Title: getContentString()}
	case fmt.Stringer:
		return &StringLink{Url: s.String(), Content: getContent(), Title: getContentString()}
	}
	v := reflect.ValueOf(url)
	if v.Kind() != reflect.String {
		panic(fmt.Errorf("Invalid type for url: %T", url))
	}
	return &StringLink{Url: v.String(), Content: getContent(), Title: getContentString()}
}

func NewLinkModelRel(url interface{}, rel string, content ...interface{}) LinkModel {
	getContent := func() View {
		if len(content) == 0 {
			return NewView(url)
		}
		return NewViews(content...)
	}
	getContentString := func() string {
		if len(content) == 0 {
			return ""
		}
		s, _ := content[0].(string)
		return s
	}
	switch s := url.(type) {
	case **Page:
		return &PageLink{Page: s, Content: getContent(), Rel: rel, Title: getContentString()}
	case *ViewWithURL:
		return &URLLink{Url: IndirectViewWithURL(s), Content: NewViews(content...), Rel: rel, Title: getContentString()}
	case LinkModel:
		if len(content) > 0 {
			return &URLLink{Url: s, Content: NewViews(content...), Rel: rel, Title: getContentString()}
		}
		return s
	case URL:
		return &URLLink{Url: s, Content: getContent(), Rel: rel, Title: getContentString()}
	case fmt.Stringer:
		return &StringLink{Url: s.String(), Content: getContent(), Rel: rel, Title: getContentString()}
	}
	v := reflect.ValueOf(url)
	if v.Kind() != reflect.String {
		panic(fmt.Errorf("Invalid type for url: %T", url))
	}
<<<<<<< HEAD
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

func (self *PageLink) URL(response *Response, args ...string) string {
	return self.Page.URL(response, args...)
}

func (self *PageLink) LinkContent(response *Response) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(response))
	}
	return self.Content
}

func (self *PageLink) LinkTitle(response *Response) string {
	if self.Title == "" {
		return self.Page.LinkTitle(response)
	}
	return self.Title
}

func (self *PageLink) LinkRel(response *Response) string {
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

func (self *StringLink) URL(response *Response, args ...string) string {
	return StringURL(self.Url).URL(response, args...)
}

func (self *StringLink) LinkContent(response *Response) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(response))
	}
	return self.Content
}

func (self *StringLink) LinkTitle(response *Response) string {
	if self.Title == "" {
		return self.Url
	}
	return self.Title
}

func (self *StringLink) LinkRel(response *Response) string {
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

func (self *URLLink) URL(response *Response, args ...string) string {
	return self.Url.URL(response, args...)
}

func (self *URLLink) LinkContent(response *Response) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(response))
	}
	return self.Content
}

func (self *URLLink) LinkTitle(response *Response) string {
	if self.Title == "" {
		return self.Url.URL(response)
	}
	return self.Title
}

func (self *URLLink) LinkRel(response *Response) string {
	return self.Rel
=======
	return &StringLink{Url: v.String(), Content: getContent(), Rel: rel, Title: getContentString()}
>>>>>>> master
}
