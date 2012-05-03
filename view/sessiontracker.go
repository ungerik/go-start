package view

import (
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// SessionTracker

type SessionTracker interface {
	ID(session *Session) (id string, ok bool)
	SetID(session *Session, id string)
	DeleteID(session *Session)
}

const sessionIdCookie = "gostart_sid"

///////////////////////////////////////////////////////////////////////////////
// CookieSessionTracker

// http://en.wikipedia.org/wiki/HTTP_cookie
type CookieSessionTracker struct {
}

func (self *CookieSessionTracker) ID(session *Session) (id string, ok bool) {
	return session.Request.GetSecureCookie(sessionIdCookie)
}

func (self *CookieSessionTracker) SetID(session *Session, id string) {
	session.Response.SetSecureCookie(sessionIdCookie, id, 0, "/")
}

func (self *CookieSessionTracker) DeleteID(session *Session) {
	session.Response.SetSecureCookie(sessionIdCookie, "delete", -time.Now().Unix(), "/")
}
