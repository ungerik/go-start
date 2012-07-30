package view

import (
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
	"net/http"
	"strconv"
	"strings"
)

func newRequest(webContext *web.Context, urlArgs []string) *Request {
	return &Request{
		Request:    webContext.Request,
		webContext: webContext,
		Params:     webContext.Params,
		URLArgs:    urlArgs,
	}
}

type Request struct {
	*http.Request
	webContext *web.Context
	Params     map[string]string
	// Arguments parsed from the URL path
	URLArgs []string
}

// URL returns the complete URL of the request including protocol and host.
func (self *Request) URLString() string {
	url := self.RequestURI
	if !strings.HasPrefix(url, "http") {
		url = "http://" + self.webContext.Request.Host + url
	}
	return url
}

// todo: all browsers
func (self *Request) ParseUserAgent() (renderer string, version utils.VersionTuple, err error) {
	s := self.UserAgent()
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

func (self *Request) Port() uint16 {
	i := strings.LastIndex(self.Host, ":")
	if i == -1 {
		return 80
	}
	port, _ := strconv.ParseInt(self.Host[i+1:], 10, 16)
	return uint16(port)
}

func (self *Request) GetSecureCookie(name string) (string, bool) {
	return self.webContext.GetSecureCookie(name)
}
