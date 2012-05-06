package view

///////////////////////////////////////////////////////////////////////////////
// Middleware

type Middleware interface {
	PreRender(response *Response) (abort bool)
	PostRender(response *Response, html string, err error) (newHtml string, newErr error)
}
