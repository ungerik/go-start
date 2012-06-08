package view

import (
	"hash/crc32"
	"bytes"
	"container/heap"
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

	dynamicStyle       dependencyHeap
	dynamicHeadScripts dependencyHeap
	dynamicScripts     dependencyHeap
}

// RequestURL returns the complete URL of the request including protocol and host.
func (self *Context) RequestURL() string {
	url := self.Request.RequestURI
	if !strings.HasPrefix(url, "http") {
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

func (self *Context) AddStyle(css string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = make(dependencyHeap, 0, 1)
		self.dynamicStyle.Init()
	}
	self.dynamicStyle.AddIfNew("<style>"+css+"</style>", priority)
}

func (self *Context) AddStyleURL(url string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = make(dependencyHeap, 0, 1)
		self.dynamicStyle.Init()
	}
	self.dynamicStyle.AddIfNew("<link rel='stylesheet' href='"+url+"'>", priority)
}

func (self *Context) AddHeaderScript(script string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = make(dependencyHeap, 0, 1)
		self.dynamicHeadScripts.Init()
	}
	self.dynamicHeadScripts.AddIfNew("<script>"+script+"</script>", priority)
}

func (self *Context) AddHeaderScriptURL(url string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = make(dependencyHeap, 0, 1)
		self.dynamicHeadScripts.Init()
	}
	self.dynamicHeadScripts.AddIfNew("<script src='"+url+"'></script>", priority)
}

func (self *Context) AddScript(script string, priority int) {
	if self.dynamicScripts == nil {
		self.dynamicScripts = make(dependencyHeap, 0, 1)
		self.dynamicScripts.Init()
	}
	self.dynamicScripts.AddIfNew("<script>"+script+"</script>", priority)
}

func (self *Context) AddScriptURL(url string, priority int) {
	if self.dynamicScripts == nil {
		self.dynamicScripts = make(dependencyHeap, 0, 1)
		self.dynamicScripts.Init()
	}
	self.dynamicScripts.AddIfNew("<script src='"+url+"'></script>", priority)
}

///////////////////////////////////////////////////////////////////////////////
// dependencyHeap

type dependencyHeapItem struct {
	text     string
	hash     uint32
	priority int
}

type dependencyHeap []dependencyHeapItem

func (self *dependencyHeap) Len() int {
	return len(*self)
}

func (self *dependencyHeap) Less(i, j int) bool {
	return (*self)[i].priority < (*self)[j].priority
}

func (self *dependencyHeap) Swap(i, j int) {
	(*self)[i], (*self)[j] = (*self)[j], (*self)[i]
}

func (self *dependencyHeap) Push(item interface{}) {
	*self = append(*self, item.(dependencyHeapItem))
}

func (self *dependencyHeap) Pop() interface{} {
	end := len(*self) - 1
	item := (*self)[end]
	*self = (*self)[:end]
	return item
}

func (self *dependencyHeap) Init() {
	heap.Init(self)
}

func (self *dependencyHeap) AddIfNew(text string, priority int) {
	hash := crc32.ChecksumIEEE([]byte(text))
	for i := range *self {
		if (*self)[i].hash == hash {
			// text is not new
			return
		}
	}
	heap.Push(self, dependencyHeapItem{text, hash, priority})
}

func (self *dependencyHeap) String() string {
	if self == nil {
		return ""
	}
	var buf bytes.Buffer
	for i := range *self {
		buf.WriteString((*self)[i].text)
	}
	return buf.String()
}
