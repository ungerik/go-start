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
	URL(request *Request, session *Session, response *Response, args ...string) string
}

///////////////////////////////////////////////////////////////////////////////
// IndirectURL

// IndirectURL encapsulates pointers to URL implementations.
// To break circular dependencies, addresses of URL implementing variables
// can be passed to this function that encapsulates it with an URL
// implementation that dereferences the pointers at runtime when they are initialized.
func IndirectURL(urlPtr interface{}) URL {
	switch s := urlPtr.(type) {
	case *URL:
		return &indirectURL{s}
	case *ViewWithURL:
		return &indirectViewWithURL{s}
	case **Page:
		return &indirectPageURL{s}
	}
	panic(errs.Format("%T not a pointer to a view.URL", urlPtr))
}

type indirectURL struct {
	url *URL
}

func (self *indirectURL) URL(request *Request, session *Session, response *Response, args ...string) string {
	return (*self.url).URL(request, session, response)
}

type indirectPageURL struct {
	page **Page
}

func (self *indirectPageURL) URL(request *Request, session *Session, response *Response, args ...string) string {
	return self.page.URL(request, session, response)
}

type indirectViewWithURL struct {
	viewWithURL *ViewWithURL
}

func (self *indirectViewWithURL) URL(request *Request, session *Session, response *Response, args ...string) string {
	return (*self.viewWithURL).URL(request, session, response)
}

///////////////////////////////////////////////////////////////////////////////
// StringURL

// StringURL implements the URL interface for a string.
type StringURL string

func (self StringURL) URL(request *Request, session *Session, response *Response, args ...string) string {
	if len(args) == 0 {
		args = request.URLArgs
	}
	url := string(self)
	for _, arg := range args {
		url = strings.Replace(url, PathFragmentPattern, arg, 1)
	}
	return url
}
