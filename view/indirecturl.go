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
	case *URL:
		return &indirectURL{s}
	case *ViewWithURL:
		return IndirectViewWithURL(s)
	case **Page:
		return &indirectPageURL{s}
	}
	panic(errs.Format("%T not a pointer to a view.URL", urlPtr))
}

type indirectURL struct {
	url *URL
}

func (self *indirectURL) URL(args ...string) string {
	return (*self.url).URL(args...)
}

type indirectPageURL struct {
	page **Page
}

func (self *indirectPageURL) URL(args ...string) string {
	return self.page.URL(args...)
}
