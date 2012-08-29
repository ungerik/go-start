package view

import (
	"github.com/ungerik/go-start/errs"
)

// IndirectURL encapsulates pointers to URL implementations.
// To break circular dependencies, addresses of URL implementing variables
// can be passed to this function that encapsulates it with an URL
// implementation that dereferences the pointers at runtime.
func IndirectURL(urlPtr interface{}) URL {
	switch s := urlPtr.(type) {
	case **Page:
		return &indirectPageURL{s}
	case *ViewWithURL:
		return IndirectViewWithURL(s)
	case *URL:
		return &indirectURL{s}
	}
	panic(errs.Format("%T not a pointer to a view.URL", urlPtr))
}

type indirectURL struct {
	url *URL
}

func (self *indirectURL) URL(ctx *Context) string {
	return (*self.url).URL(ctx)
}

type indirectPageURL struct {
	page **Page
}

func (self *indirectPageURL) URL(ctx *Context) string {
	return self.page.URL(ctx)
}
