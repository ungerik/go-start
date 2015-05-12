package view

import (
	// "encoding/base64"
	"log"
	"sync"
	"time"

	"github.com/ungerik/go-start/utils"
)

///////////////////////////////////////////////////////////////////////////////
// SessionUserAuth

// SessionUserAuth implements Authenticator by
type SessionUserAuth struct {
	LoginURL     URL
	UserPassword map[string]string // map of username to base64 encoded SHA1 hash of the password
	userLogins   map[string]time.Time
	mutex        sync.Mutex
}

func NewSessionUserAuth(usernamesAndPasswords ...string) *SessionUserAuth {
	userPass := make(map[string]string, len(usernamesAndPasswords)/2)

	for i := 0; i < len(usernamesAndPasswords)/2; i++ {
		username := usernamesAndPasswords[i*2]
		password := usernamesAndPasswords[i*2+1]
		userPass[username] = utils.SHA1Base64String(password)
	}

	return &SessionUserAuth{
		UserPassword: userPass,
		userLogins:   make(map[string]time.Time),
	}
}

func (auth *SessionUserAuth) Authenticate(ctx *Context) (ok bool, err error) {
	auth.mutex.Lock()
	defer auth.mutex.Unlock()

	user := ctx.Session.ID()

	_, hasUser := auth.UserPassword[user]
	if !hasUser || user == "" {
		auth.logout(user, ctx)
		return false, nil
	}

	userLoginTime, hasLogin := auth.userLogins[user]
	if !hasLogin || time.Now().After(userLoginTime.Add(Config.SessionUserAuthLogoutAfter)) {
		auth.logout(user, ctx)
		return false, nil
	}

	ctx.AuthUser = user
	return true, nil
}

func (auth *SessionUserAuth) Login(username, password string, ctx *Context) (ok bool, err error) {
	auth.mutex.Lock()
	defer auth.mutex.Unlock()

	passHash, hasUser := auth.UserPassword[username]
	if hasUser && passHash == utils.SHA1Base64String(password) {
		ctx.Session.SetID(username)
		ctx.AuthUser = username
		auth.userLogins[username] = time.Now()
		log.Println("Logged in user", username)
		return true, nil
	}

	ctx.Session.DeleteID()
	ctx.AuthUser = ""
	return false, nil
}

func (auth *SessionUserAuth) logout(user string, ctx *Context) {
	delete(auth.userLogins, user)
	ctx.AuthUser = ""
	ctx.Session.DeleteID()
	if auth.LoginURL != nil {
		ctx.Response.RedirectTemporary302(auth.LoginURL.URL(ctx))
	}
}

func (auth *SessionUserAuth) Logout(ctx *Context) {
	user := ctx.Session.ID()
	if user == "" {
		return
	}

	auth.mutex.Lock()
	defer auth.mutex.Unlock()

	auth.logout(user, ctx)
	log.Println("Logged out user", user)
}

///////////////////////////////////////////////////////////////////////////////
// HtpasswdWatchingSessionUserAuth

type HtpasswdWatchingSessionUserAuth struct {
	SessionUserAuth
	htpasswdFile     string
	htpasswdFileTime time.Time
}

func NewHtpasswdWatchingSessionUserAuth(htpasswdFile string) (auth *HtpasswdWatchingSessionUserAuth, err error) {
	auth = &HtpasswdWatchingSessionUserAuth{
		SessionUserAuth: SessionUserAuth{userLogins: make(map[string]time.Time)},
		htpasswdFile:    htpasswdFile,
	}
	auth.UserPassword, auth.htpasswdFileTime, err = utils.ReadHtpasswdFile(htpasswdFile)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (auth *HtpasswdWatchingSessionUserAuth) Authenticate(ctx *Context) (ok bool, err error) {
	auth.mutex.Lock()

	t, err := utils.FileModifiedTime(auth.htpasswdFile)
	if err != nil {
		auth.mutex.Unlock()
		return false, err
	}

	if t.After(auth.htpasswdFileTime) {
		auth.UserPassword, auth.htpasswdFileTime, err = utils.ReadHtpasswdFile(auth.htpasswdFile)
		if err != nil {
			auth.mutex.Unlock()
			return false, err
		}
	}

	auth.mutex.Unlock()

	return auth.SessionUserAuth.Authenticate(ctx)
}
