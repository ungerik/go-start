package view

import (
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// SessionTracker

type SessionTracker interface {
	ID(ctx *Context) (id string)
	SetID(ctx *Context, id string)
	DeleteID(ctx *Context)
}

const sessionIdCookie = "gostart_sid"

///////////////////////////////////////////////////////////////////////////////
// CookieSessionTracker

// http://en.wikipedia.org/wiki/HTTP_cookie
type CookieSessionTracker struct {
}

func (self *CookieSessionTracker) ID(ctx *Context) string {
	id, _ := ctx.Request.GetSecureCookie(sessionIdCookie)
	return id
}

func (self *CookieSessionTracker) SetID(ctx *Context, id string) {
	ctx.Response.SetSecureCookie(sessionIdCookie, id, 0, "/")
}

func (self *CookieSessionTracker) DeleteID(ctx *Context) {
	ctx.Response.SetSecureCookie(sessionIdCookie, "delete", -time.Now().Unix(), "/")
}
