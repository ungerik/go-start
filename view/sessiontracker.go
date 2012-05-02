package view

import (
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// SessionTracker

type SessionTracker interface {
	ID(request *Request, session *Session, response *Response) (id string, ok bool)
	SetID(request *Request, session *Session, response *Response, id string)
	DeleteID(request *Request, session *Session, response *Response)
}

const sessionIdCookie = "gostart_sid"

///////////////////////////////////////////////////////////////////////////////
// CookieSessionTracker

// http://en.wikipedia.org/wiki/HTTP_cookie
type CookieSessionTracker struct {
}

func (self *CookieSessionTracker) ID(request *Request, session *Session, response *Response) (id string, ok bool) {
	return context.GetSecureCookie(sessionIdCookie)
}

func (self *CookieSessionTracker) SetID(request *Request, session *Session, response *Response, id string) {
	context.SetSecureCookie(sessionIdCookie, id, 0, "/")
}

func (self *CookieSessionTracker) DeleteID(request *Request, session *Session, response *Response) {
	context.SetSecureCookie(sessionIdCookie, "delete", -time.Now().Unix(), "/")
}
