package view

import (
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// SessionTracker

type SessionTracker interface {
	ID(context *Context) (id string, ok bool)
	SetID(context *Context, id string)
	DeleteID(context *Context)
}

const sessionIdCookie = "gostart_sid"

///////////////////////////////////////////////////////////////////////////////
// CookieSessionTracker

// http://en.wikipedia.org/wiki/HTTP_cookie
type CookieSessionTracker struct {
}

func (self *CookieSessionTracker) ID(context *Context) (id string, ok bool) {
	return context.GetSecureCookie(sessionIdCookie)
}

func (self *CookieSessionTracker) SetID(context *Context, id string) {
	context.SetSecureCookie(sessionIdCookie, id, 0, "/")
}

func (self *CookieSessionTracker) DeleteID(context *Context) {
	context.SetSecureCookie(sessionIdCookie, "delete", -time.Now().Unix(), "/")
}
