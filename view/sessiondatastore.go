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
	Get(session *Session, data interface{}) (ok bool, err error)
	Set(session *Session, data interface{}) (err error)
	Delete(session *Session) (err error)
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

func (self *CookieSessionDataStore) Get(session *Session, data interface{}) (ok bool, err error) {
	sessionID, ok := session.ID()
	if !ok {
		return false, errs.Format("Can't set session data without a session id")
	}

	cookieValue, ok := session.Request.GetSecureCookie(self.cookieName(sessionID))
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

func (self *CookieSessionDataStore) Set(session *Session, data interface{}) (err error) {
	sessionID, ok := session.ID()
	if !ok {
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

	session.Response.SetSecureCookie(self.cookieName(sessionID), string(dataBytes), 0, "/")
	return nil
}

func (self *CookieSessionDataStore) Delete(session *Session) (err error) {
	sessionID, ok := session.ID()
	if !ok {
		return errs.Format("Can't delete session data without a session id")
	}

	session.Response.SetSecureCookie(self.cookieName(sessionID), "", -time.Now().Unix(), "/")
	return nil
}
