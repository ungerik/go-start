package view

import (
	"bytes"
	"container/heap"
	"fmt"
	"hash/crc32"
	"net/http"

	"github.com/ungerik/web.go"
)

func newResponse(webContext *web.Context, respondingView View, urlArgs []string) *Response {
	response := &Response{
		webContext:     webContext,
		RespondingView: respondingView,
		Request:        newRequest(webContext, urlArgs),
		Session:        new(Session),
	}
	response.Session.init(response.Request, response)
	return response
}

type Response struct {
	buffer     bytes.Buffer
	webContext *web.Context

	Request *Request
	Session *Session

	// View that responds to the HTTP request
	RespondingView View
	// Custom response wide data that can be set by the application
	Data interface{}

	dynamicStyle       dependencyHeap
	dynamicHeadScripts dependencyHeap
	dynamicScripts     dependencyHeap
}

// New creates a clone of the response with an empty buffer.
// Used to render preliminary text.
func (self *Response) New() *Response {
	return &Response{
		webContext:     self.webContext,
		Request:        self.Request,
		Session:        self.Session,
		RespondingView: self.RespondingView,
		Data:           self.Data,
	}
}

func (self *Response) Write(p []byte) (n int, err error) {
	return self.buffer.Write(p)
}

func (self *Response) WriteByte(c byte) error {
	return self.buffer.WriteByte(c)
}

func (self *Response) WriteRune(r rune) (n int, err error) {
	return self.buffer.WriteRune(r)
}

func (self *Response) WriteString(s string) (n int, err error) {
	return self.buffer.WriteString(s)
}

func (self *Response) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(&self.buffer, format, args...)
}

func (self *Response) String() string {
	return self.buffer.String()
}

func (self *Response) SetSecureCookie(name string, val string, age int64, path string) {
	self.webContext.SetSecureCookie(name, val, age, path)
}

func (self *Response) Abort(status int, body string) {
	self.webContext.Abort(status, body)
}

func (self *Response) RedirectPermanently301(url string) {
	self.webContext.Redirect(301, url)
}

func (self *Response) RedirectTemporary302(url string) {
	self.webContext.Redirect(302, url)
}

func (self *Response) NotModified304() {
	self.webContext.NotModified()
}

func (self *Response) Forbidden403(message string) {
	self.Abort(403, message)
}

func (self *Response) NotFound404(message string) {
	self.Abort(404, message)
}

func (self *Response) AuthorizationRequired401() {
	self.Abort(401, "Authorization Required")
}

func (self *Response) Header() http.Header {
	return self.webContext.Header()
}

func (self *Response) ContentType(ext string) {
	self.webContext.ContentType(ext)
}

func (self *Response) AddStyle(css string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = make(dependencyHeap, 0, 1)
		self.dynamicStyle.Init()
	}
	self.dynamicStyle.AddIfNew("<style>"+css+"</style>", priority)
}

func (self *Response) AddStyleURL(url string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = make(dependencyHeap, 0, 1)
		self.dynamicStyle.Init()
	}
	self.dynamicStyle.AddIfNew("<link rel='stylesheet' href='"+url+"'>", priority)
}

func (self *Response) AddHeaderScript(script string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = make(dependencyHeap, 0, 1)
		self.dynamicHeadScripts.Init()
	}
	self.dynamicHeadScripts.AddIfNew("<script>"+script+"</script>", priority)
}

func (self *Response) AddHeaderScriptURL(url string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = make(dependencyHeap, 0, 1)
		self.dynamicHeadScripts.Init()
	}
	self.dynamicHeadScripts.AddIfNew("<script src='"+url+"'></script>", priority)
}

func (self *Response) AddScript(script string, priority int) {
	if self.dynamicScripts == nil {
		self.dynamicScripts = make(dependencyHeap, 0, 1)
		self.dynamicScripts.Init()
	}
	self.dynamicScripts.AddIfNew("<script>"+script+"</script>", priority)
}

func (self *Response) AddScriptURL(url string, priority int) {
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
