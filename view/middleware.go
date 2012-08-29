package view

///////////////////////////////////////////////////////////////////////////////
// Middleware

type Middleware interface {
	PreRender(ctx *Context) (abort bool)
	PostRender(response *Response, html string, err error) (newHtml string, newErr error)
}
