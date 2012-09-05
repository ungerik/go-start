package view

import (
	"fmt"

	"github.com/ungerik/web.go"
)

func newContext(webContext *web.Context, respondingView View, urlArgs []string) *Context {
	ctx := &Context{
		URLArgs:  urlArgs,
		Request:  newRequest(webContext),
		Response: newResponse(webContext, respondingView),
	}
	ctx.Session = newSession(ctx)
	return ctx
}

type Context struct {
	Request  *Request
	Response *Response
	Session  *Session

	// Arguments parsed from the URL path
	URLArgs []string

	// Custom response wide data that can be set by the application
	Data      interface{}
	DebugData interface{}
}

/*
ForURLArgs returns an altered Context copy where
Context.URLArgs is set to urlArgs.
Can be used for calling the the URL() method of a URL interface
to get the URL of another view, defined by urlArgs.

The following example gets the URL of MyPage with the first
URL argument is that of the current page and the second
URL argument is "second-arg":

	MyPage.URL(ctx.ForURLArgs(ctx.URLArgs[0], "second-arg"))
*/
func (self *Context) ForURLArgs(urlArgs ...string) *Context {
	clone := *self
	clone.URLArgs = urlArgs
	return &clone
}

func (self *Context) ForURLArgsConvert(urlArgs ...interface{}) *Context {
	stringArgs := make([]string, len(urlArgs))
	for i := range urlArgs {
		stringArgs[i] = fmt.Sprint(urlArgs[i])
	}
	return self.ForURLArgs(stringArgs...)
}
