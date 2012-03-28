package view

import (
	"github.com/ungerik/go-start/errs"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// URL

// URL is an interface to return URL strings depending on the request context.
type URL interface {
	// If args are passed, they will be used instead of context.PathArgs.
	URL(context *Context, args ...string) string
}

// IndirectURL encapsulates pointers to URL implementations.
// To break circular dependencies, addresses of URL implementing variables
// can be passed to this function that encapsulates it with an URL
// implementation that dereferences the pointers at runtime when they are initialized.
func IndirectURL(urlPtr interface{}) URL {
	switch s := urlPtr.(type) {
	case *URL:
		return &indirectURL{s}
	case *ViewWithURL:
		return &IndirectViewWithURL{s}
	case **Page:
		return &IndirectPageURL{s}
	}
	panic(errs.Format("%T not a pointer to a view.URL", urlPtr))
}

type indirectURL struct {
	url *URL
}

func (self *indirectURL) URL(context *Context, args ...string) string {
	return (*self.url).URL(context)
}

///////////////////////////////////////////////////////////////////////////////
// StringURL

// StringURL implements the URL interface for a string.
type StringURL string

func (self StringURL) URL(context *Context, args ...string) string {
	if len(args) == 0 {
		args = context.PathArgs
	}
	url := string(self)
	for _, arg := range args {
		url = strings.Replace(url, PathFragmentPattern, arg, 1)
	}
	return url
}

///////////////////////////////////////////////////////////////////////////////
// PageURL

// IndirectPageURL implements the URL interface for pointer to a Page pointer.
// Used to break circular dependencies by referencing a Page pointer variable.
type IndirectPageURL struct {
	Page **Page
}

func (self *IndirectPageURL) URL(context *Context, args ...string) string {
	return self.Page.URL(context)
}

///////////////////////////////////////////////////////////////////////////////
// IndirectViewWithURL

// IndirectViewWithURL implements the URL interface for pointer to a ViewWithURL.
// Used to break circular dependencies by referencing a ViewWithURL variable.
type IndirectViewWithURL struct {
	ViewWithURL *ViewWithURL
}

func (self *IndirectViewWithURL) URL(context *Context, args ...string) string {
	return (*self.ViewWithURL).URL(context)
}
