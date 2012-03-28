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
	Get(context *Context, data interface{}) (ok bool, err error)
	Set(context *Context, data interface{}) (err error)
	Delete(context *Context) (err error)
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

func (self *CookieSessionDataStore) Get(context *Context, data interface{}) (ok bool, err error) {
	sessionID, ok := context.SessionID()
	if !ok {
		return false, errs.Format("Can't set session data without a session id")
	}

	cookieValue, ok := context.GetSecureCookie(self.cookieName(sessionID))
	if !ok {
		return false, nil
	}

	decryptedCookieValue, err := context.DecryptCookie([]byte(cookieValue))
	if err != nil {
		return false, err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(decryptedCookieValue))
	err = decoder.Decode(data)
	return err == nil, err
}

func (self *CookieSessionDataStore) Set(context *Context, data interface{}) (err error) {
	sessionID, ok := context.SessionID()
	if !ok {
		return errs.Format("Can't set session data without a session id")
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	dataBytes, err := context.EncryptCookie(buffer.Bytes())
	if err != nil {
		return err
	}

	if len(dataBytes) > 4000 { // Max HTTP header size is 4094 minus some space for protocol
		return errs.Format("Session %s data size %d is larger than cookie limit of 4000 bytes", sessionID, len(dataBytes))
	}

	context.SetSecureCookie(self.cookieName(sessionID), string(dataBytes), 0, "/")
	return nil
}

func (self *CookieSessionDataStore) Delete(context *Context) (err error) {
	sessionID, ok := context.SessionID()
	if !ok {
		return errs.Format("Can't delete session data without a session id")
	}

	context.SetSecureCookie(self.cookieName(sessionID), "", -time.Now().Unix(), "/")
	return nil
}
