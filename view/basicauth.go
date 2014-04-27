package view

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

func stringSHA1Base64(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// NewBasicAuth creates a BasicAuth instance with a single username and password.
func NewBasicAuth(realm, username, password string) *BasicAuth {
	return &BasicAuth{
		Realm:        realm,
		UserPassword: map[string]string{username: stringSHA1Base64(password)},
	}
}

func NewBasicAuthFromHtpasswd(realm, filename string) (*BasicAuth, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ':'
	csvReader.Comment = '#'
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	basicAuth := &BasicAuth{
		Realm:        realm,
		UserPassword: make(map[string]string),
	}

	for _, record := range records {
		username := record[0]
		password := record[1]
		if len(password) < 5 {
			return nil, errors.New("Invalid password")
		}
		if password[:5] != "{SHA}" {
			return nil, errors.New("Unsupported password format, must be SHA1")
		}
		basicAuth.UserPassword[username] = password[5:]
	}

	return basicAuth, nil
}

///////////////////////////////////////////////////////////////////////////////
// BasicAuth

// BasicAuth implements HTTP basic auth as Authenticator.
type BasicAuth struct {
	Realm        string
	UserPassword map[string]string // map of username to base64 encoded SHA1 hash of the password
}

func (self *BasicAuth) Authenticate(ctx *Context) (ok bool, err error) {
	header := ctx.Request.Header.Get("Authorization")
	f := strings.Fields(header)
	if len(f) == 2 && f[0] == "Basic" {
		if b, err := base64.StdEncoding.DecodeString(f[1]); err == nil {
			a := strings.Split(string(b), ":")
			if len(a) == 2 {
				username := a[0]
				password := a[1]
				p, ok := self.UserPassword[username]
				if ok && p == stringSHA1Base64(password) {
					return true, nil
				}
			}
		}
	}

	ctx.Response.Header().Set("WWW-Authenticate", "Basic realm=\""+self.Realm+"\"")
	ctx.Response.AuthorizationRequired401()
	return false, nil
}
