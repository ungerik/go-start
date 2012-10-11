package view

import (
	"bytes"
	"encoding/gob"
	"github.com/ungerik/go-start/errs"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// SessionDataStore

type SessionDataStore interface {
	Get(ctx *Context, data interface{}) (ok bool, err error)
	Set(ctx *Context, data interface{}) (err error)
	Delete(ctx *Context) (err error)
}

///////////////////////////////////////////////////////////////////////////////
// CookieSessionDataStore

func NewCookieSessionDataStore() SessionDataStore {
	return &CookieSessionDataStore{"gostart_session_data"}
}

type CookieSessionDataStore struct {
	cookieNameBase string
}

func (self *CookieSessionDataStore) cookieName(sessionID string) string {
	return self.cookieNameBase + sessionID
}

func (self *CookieSessionDataStore) Get(ctx *Context, data interface{}) (ok bool, err error) {
	sessionID := ctx.Session.ID()
	if sessionID == "" {
		return false, errs.Format("Can't set session data without a session id")
	}

	cookieValue, ok := ctx.Request.GetSecureCookie(self.cookieName(sessionID))
	if !ok {
		return false, nil
	}

	decryptedCookieValue, err := DecryptCookie([]byte(cookieValue))
	if err != nil {
		return false, err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(decryptedCookieValue))
	err = decoder.Decode(data)
	return err == nil, err
}

func (self *CookieSessionDataStore) Set(ctx *Context, data interface{}) (err error) {
	sessionID := ctx.Session.ID()
	if sessionID == "" {
		return errs.Format("Can't set session data without a session id")
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	dataBytes, err := EncryptCookie(buffer.Bytes())
	if err != nil {
		return err
	}

	if len(dataBytes) > 4000 { // Max HTTP header size is 4094 minus some space for protocol
		return errs.Format("Session %s data size %d is larger than cookie limit of 4000 bytes", sessionID, len(dataBytes))
	}

	ctx.Response.SetSecureCookie(self.cookieName(sessionID), string(dataBytes), 0, "/")
	return nil
}

func (self *CookieSessionDataStore) Delete(ctx *Context) (err error) {
	sessionID := ctx.Session.ID()
	if sessionID == "" {
		return errs.Format("Can't delete session data without a session id")
	}

	ctx.Response.SetSecureCookie(self.cookieName(sessionID), "", -time.Now().Unix(), "/")
	return nil
}
