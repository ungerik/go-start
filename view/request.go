package view

import (
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
	"net/http"
	"strconv"
	"strings"
)

func newRequest(webContext *web.Context) *Request {
	return &Request{
		Request:    webContext.Request,
		webContext: webContext,
		Params:     webContext.Params,
	}
}

type Request struct {
	*http.Request
	webContext *web.Context
	Params     map[string]string
}

// AddProtocolAndHostToURL adds the protocol (http:// or https://)
// and request host (domain or IP) to an URL if not present.
func (self *Request) AddProtocolAndHostToURL(url string) string {
	if len(url) > 0 && url[0] == '/' {
		if self.TLS != nil {
			url = "https://" + self.Host + url
		} else {
			url = "http://" + self.Host + url
		}
	}
	return url
}

// URL returns the complete URL of the request including protocol and host.
func (self *Request) URLString() string {
	return self.AddProtocolAndHostToURL(self.RequestURI)
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
