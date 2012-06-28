package view

import (
	"fmt"
	"reflect"
)

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
		return &URLLink{Url: IndirectViewWithURL(s), Content: NewViews(content...)}
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
		return &URLLink{Url: IndirectViewWithURL(s), Content: NewViews(content...), Rel: rel}
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
