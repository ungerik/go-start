package view

import (
	"encoding/base64"
	"strings"
	"sync"
	"time"

	"github.com/ungerik/go-start/utils"
)

// BasicAuth implements HTTP basic auth as Authenticator.
// See also HtpasswdWatchingBasicAuth
type BasicAuth struct {
	Realm        string
	UserPassword map[string]string // map of username to base64 encoded SHA1 hash of the password
}

// NewBasicAuth creates a BasicAuth instance with a series of
// not encrypted usernames and passwords that are provided in alternating order.
func NewBasicAuth(realm string, usernamesAndPasswords ...string) *BasicAuth {
	userPass := make(map[string]string, len(usernamesAndPasswords)/2)

	for i := 0; i < len(usernamesAndPasswords)/2; i++ {
		username := usernamesAndPasswords[i*2]
		password := usernamesAndPasswords[i*2+1]
		userPass[username] = utils.SHA1Base64String(password)
	}

	return &BasicAuth{
		Realm:        realm,
		UserPassword: userPass,
	}
}

func (basicAuth *BasicAuth) Authenticate(ctx *Context) (ok bool, err error) {
	header := ctx.Request.Header.Get("Authorization")
	f := strings.Fields(header)
	if len(f) == 2 && f[0] == "Basic" {
		if b, err := base64.StdEncoding.DecodeString(f[1]); err == nil {
			a := strings.Split(string(b), ":")
			if len(a) == 2 {
				username := a[0]
				password := a[1]
				p, ok := basicAuth.UserPassword[username]
				if ok && p == utils.SHA1Base64String(password) {
					ctx.AuthUser = username
					return true, nil
				}
			}
		}
	}

	ctx.Response.Header().Set("WWW-Authenticate", "Basic realm=\""+basicAuth.Realm+"\"")
	ctx.Response.AuthorizationRequired401()
	return false, nil
}

type HtpasswdWatchingBasicAuth struct {
	basicAuth        BasicAuth
	htpasswdFile     string
	htpasswdFileTime time.Time
	mutex            sync.Mutex
}

func NewHtpasswdWatchingBasicAuth(realm, htpasswdFile string) (auth *HtpasswdWatchingBasicAuth, err error) {
	auth = &HtpasswdWatchingBasicAuth{
		basicAuth:    BasicAuth{Realm: realm},
		htpasswdFile: htpasswdFile,
	}
	auth.basicAuth.UserPassword, auth.htpasswdFileTime, err = utils.ReadHtpasswdFile(htpasswdFile)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (auth *HtpasswdWatchingBasicAuth) Authenticate(ctx *Context) (ok bool, err error) {
	auth.mutex.Lock()
	defer auth.mutex.Unlock()

	t, err := utils.FileModifiedTime(auth.htpasswdFile)
	if err != nil {
		return false, err
	}

	if t.After(auth.htpasswdFileTime) {
		auth.basicAuth.UserPassword, auth.htpasswdFileTime, err = utils.ReadHtpasswdFile(auth.htpasswdFile)
		if err != nil {
			return false, err
		}
	}

	return auth.basicAuth.Authenticate(ctx)
}
