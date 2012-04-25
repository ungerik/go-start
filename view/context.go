package view

import (
	"encoding/base64"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
	"strconv"
	"strings"
)

func NewContext(webContext *web.Context, respondingView View, pathArgs []string) *Context {
	return &Context{
		web.Context:    webContext,
		RespondingView: respondingView,
		PathArgs:       pathArgs,
	}
}

///////////////////////////////////////////////////////////////////////////////
// Context

// Context holds all data specific to a HTTP request and will be passed to View.Render() methods.
type Context struct {
	*web.Context

	// View that responds to the HTTP request
	RespondingView View

	// Arguments parsed from the URL path
	PathArgs []string

	/*
		Cached user object of the session.
		User won't be set automatically, use user.OfSession(context) instead.

		Example for setting it automatically for every request:

			import "github.com/ungerik/go-start/user"

			Config.OnPreAuth = func(context *Context) error {
				user.OfSession(context) // Sets context.User
				return nil
			}
	*/
	User interface{}

	// Custom request wide data that can be set by the application
	Data interface{}

	cachedSessionID string
	//	cache           map[string]interface{}
}

// RequestURL returns the complete URL of the request including protocol and host.
func (self *Context) RequestURL() string {
	url := self.Request.RequestURI
	if !utils.StringStartsWith(url, "http") {
		url = "http://" + self.Request.Host + url
	}
	return url
}

func (self *Context) EncryptCookie(data []byte) (result []byte, err error) {
	// todo crypt

	e := base64.StdEncoding
	result = make([]byte, e.EncodedLen(len(data)))
	e.Encode(result, data)
	return result, nil
}

func (self *Context) DecryptCookie(data []byte) (result []byte, err error) {
	// todo crypt

	e := base64.StdEncoding
	result = make([]byte, e.DecodedLen(len(data)))
	_, err = e.Decode(result, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// todo: all browsers
func (self *Context) ParseRequestUserAgent() (renderer string, version utils.VersionTuple, err error) {
	s := self.Request.UserAgent()
	switch {
	case strings.Contains(s, "Gecko"):
		if i := strings.Index(s, "rv:"); i != -1 {
			i += len("rv:")
			if l := strings.IndexAny(s[i:], "); "); l != -1 {
				if version, err = utils.ParseVersionTuple(s[i : i+l]); err == nil {
					return "Gecko", version, nil
				}
			}
		}
	case strings.Contains(s, "MSIE "):
		if i := strings.Index(s, "MSIE "); i != -1 {
			i += len("MSIE ")
			if l := strings.IndexAny(s[i:], "); "); l != -1 {
				if version, err = utils.ParseVersionTuple(s[i : i+l]); err == nil {
					return "MSIE", version, nil
				}
			}
		}
	}
	return "", nil, nil
}

func (self *Context) RequestPort() uint16 {
	i := strings.LastIndex(self.Request.Host, ":")
	if i == -1 {
		return 80
	}
	port, _ := strconv.ParseInt(self.Request.Host[i+1:], 10, 16)
	return uint16(port)
}

//func (self *Context) Cache(key string, value interface{}) {
//	if self.cache == nil {
//		self.cache = make(map[string]interface{})
//	}
//	self.cache[key] = value
//}
//
//func (self *Context) Cached(key string) (value interface{}, ok bool) {
//	if self.cache == nil {
//		return nil, false
//	}
//	value, ok = self.cache[key]
//	return value, ok
//}
//
//func (self *Context) DeleteCached(key string) {
//	if self.cache == nil {
//		return
//	}
//	self.cache[key] = nil, false
//}

// SessionID returns the id of the session and if there is a session active.
func (self *Context) SessionID() (id string, ok bool) {
	if self.cachedSessionID != "" {
		return self.cachedSessionID, true
	}

	if Config.SessionTracker == nil {
		return "", false
	}

	self.cachedSessionID, ok = Config.SessionTracker.ID(self)
	return self.cachedSessionID, ok
}

func (self *Context) SetSessionID(id string) {
	if Config.SessionTracker != nil {
		Config.SessionTracker.SetID(self, id)
		self.cachedSessionID = id
	}
}

func (self *Context) DeleteSessionID() {
	self.cachedSessionID = ""
	if t := Config.SessionTracker; t != nil {
		t.DeleteID(self)
	}
}

// SessionData returns all session data in out.
func (self *Context) SessionData(out interface{}) (ok bool, err error) {
	if Config.SessionDataStore == nil {
		return false, errs.Format("Can't get session data without gostart/views.Config.SessionDataStore")
	}
	return Config.SessionDataStore.Get(self, out)
}

// SetSessionData sets all session data.
func (self *Context) SetSessionData(data interface{}) (err error) {
	if Config.SessionDataStore == nil {
		return errs.Format("Can't set session data without gostart/views.Config.SessionDataStore")
	}
	return Config.SessionDataStore.Set(self, data)
}

// DeleteSessionData deletes all session data.
func (self *Context) DeleteSessionData() (err error) {
	if Config.SessionDataStore == nil {
		return errs.Format("Can't delete session data without gostart/views.Config.SessionDataStore")
	}
	return Config.SessionDataStore.Delete(self)
}
