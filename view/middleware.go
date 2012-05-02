package view

///////////////////////////////////////////////////////////////////////////////
// Middleware

type Middleware interface {
	PreRender(request *Request, session *Session, response *Response) (abort bool)
	PostRender(request *Request, session *Session, response *Response, html string, err error) (newHtml string, newErr error)
}
