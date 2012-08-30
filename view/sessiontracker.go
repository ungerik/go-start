package view

import (
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// SessionTracker

type SessionTracker interface {
	ID(ctx *Context) (id string, ok bool)
	SetID(ctx *Context, id string)
	DeleteID(ctx *Context)
}

const sessionIdCookie = "gostart_sid"

///////////////////////////////////////////////////////////////////////////////
// CookieSessionTracker

// http://en.wikipedia.org/wiki/HTTP_cookie
type CookieSessionTracker struct {
}

func (self *CookieSessionTracker) ID(ctx *Context) (id string, ok bool) {
	return ctx.Request.GetSecureCookie(sessionIdCookie)
}

func (self *CookieSessionTracker) SetID(ctx *Context, id string) {
	ctx.Response.SetSecureCookie(sessionIdCookie, id, 0, "/")
}

func (self *CookieSessionTracker) DeleteID(ctx *Context) {
	ctx.Response.SetSecureCookie(sessionIdCookie, "delete", -time.Now().Unix(), "/")
}
