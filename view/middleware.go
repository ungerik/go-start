package view

///////////////////////////////////////////////////////////////////////////////
// Middleware

type Middleware interface {
	PreRender(context *Context) (abort bool)
	PostRender(context *Context, html string, err error) (newHtml string, newErr error)
}
