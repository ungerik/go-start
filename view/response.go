package view

import (
	"bytes"
	"container/heap"
	"fmt"
	"hash/crc32"
	"net/http"

	"github.com/ungerik/go-start/utils"
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
	response.PushBuffer()
	return response
}

type responseBuffer struct {
	buffer *bytes.Buffer
	xml    *utils.XMLWriter
}

type Response struct {
	webContext *web.Context

	Request *Request
	Session *Session

	// View that responds to the HTTP request
	RespondingView View

	bufferStack []responseBuffer
	// XML allowes the Response to be used as utils.XMLWriter
	XML *utils.XMLWriter

	// Custom response wide data that can be set by the application
	Data      interface{}
	DebugData interface{}

	dynamicStyle       dependencyHeap
	dynamicHeadScripts dependencyHeap
	dynamicScripts     dependencyHeap
}

// PushBuffer pushes the buffer of the response body on a stack
// and sets a new empty buffer.
// This can be used to render intermediate text results.
// Note: Only the response body is pushed, all other state changes
// like setting headers will affect the final response.
func (self *Response) PushBuffer() {
	var b responseBuffer
	b.buffer = new(bytes.Buffer)
	b.xml = utils.NewXMLWriter(b.buffer)
	self.bufferStack = append(self.bufferStack, b)
	self.XML = b.xml
}

// PopBuffer pops the buffer of the response body from the stack
// and returns its content.
// This can be used to render intermediate text results.
func (self *Response) PopBuffer() (bufferData []byte) {
	last := len(self.bufferStack) - 1
	bufferData = self.bufferStack[last].buffer.Bytes()
	self.bufferStack = self.bufferStack[0:last]
	self.XML = self.bufferStack[last-1].xml
	return bufferData
}

// PopBufferString pops the buffer of the response body from the stack
// and returns its content as string.
// This can be used to render intermediate text results.
func (self *Response) PopBufferString() (bufferData string) {
	return string(self.PopBuffer())
}

/*
URLArgs returns an altered Response copy where
Response.Request.URLArgs are set to urlArgs.
Used when calling the the URL() method of a URL interface
to get the URL of another view, defined by urlArgs.

The following example gets the URL of MyPage with the first
URL argument is that of the current page and the second
URL argument is "second-arg":

	MyPage.URL(response.URLArgs(response.Request.URLArgs[0], "second-arg"))
*/
func (self *Response) URLArgs(urlArgs ...string) *Response {
	copy := *self
	copy.Request = self.Request.cloneWithURLArgs(urlArgs)
	return &copy
}

func (self *Response) Write(p []byte) (n int, err error) {
	return self.XML.Write(p)
}

func (self *Response) WriteByte(c byte) error {
	_, err := self.XML.Write([]byte{c})
	return err
}

func (self *Response) WriteString(s string) (n int, err error) {
	return self.XML.Write([]byte(s))
}

func (self *Response) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(self.XML, format, args...)
}

func (self *Response) String() string {
	return self.bufferStack[len(self.bufferStack)-1].buffer.String()
}

func (self *Response) Bytes() []byte {
	return self.bufferStack[len(self.bufferStack)-1].buffer.Bytes()
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

func (self *Response) ContentTypeByExt(ext string) {
	self.webContext.ContentType(ext)
}

// ContentDispositionAttachment makes the webbrowser open a
// "Save As.." dialog for the response.
func (self *Response) ContentDispositionAttachment(filename string) {
	self.Header().Add("Content-Type", "application/x-unknown")
	self.Header().Add("Content-Disposition", "attachment;filename="+filename)
}

// AddStyle adds dynamic CSS content to the page.
// Multiple dynamic entries will be sorted by priority.
// Dynamic CSS will be inserted after the regular CSS of the page.
// 
// Use this feature to dynamically add CSS to the page if the
// HTML content requires it.
func (self *Response) AddStyle(css string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = newDependencyHeap()
	}
	self.dynamicStyle.AddIfNew("<style>"+css+"</style>", priority)
}

// AddStyleURL adds a dynamic CSS link to the page.
// Multiple dynamic entries will be sorted by priority.
// Dynamic CSS will be inserted after the regular CSS of the page.
// 
// Use this feature to dynamically add CSS to the page if the
// HTML content requires it.
func (self *Response) AddStyleURL(url string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = newDependencyHeap()
	}
	self.dynamicStyle.AddIfNew("<link rel='stylesheet' href='"+url+"'>", priority)
}

// AddHeaderScript adds dynamic JavaScript to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// head-scripts of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) AddHeaderScript(script string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = newDependencyHeap()
	}
	self.dynamicHeadScripts.AddIfNew("<script>"+script+"</script>", priority)
}

// AddHeaderScriptURL adds a dynamic JavaScript link to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// head-scripts of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) AddHeaderScriptURL(url string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = newDependencyHeap()
	}
	self.dynamicHeadScripts.AddIfNew("<script src='"+url+"'></script>", priority)
}

// AddScript adds dynamic JavaScript to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// scripts near the end of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) AddScript(script string, priority int) {
	if self.dynamicScripts == nil {
		self.dynamicScripts = newDependencyHeap()
	}
	self.dynamicScripts.AddIfNew("<script>"+script+"</script>", priority)
}

// AddScriptURL adds a dynamic JavaScript link to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// scripts near the end of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) AddScriptURL(url string, priority int) {
	if self.dynamicScripts == nil {
		self.dynamicScripts = newDependencyHeap()
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

func newDependencyHeap() dependencyHeap {
	dh := make(dependencyHeap, 0, 1)
	heap.Init(&dh)
	return dh
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

func (self *dependencyHeap) AddIfNew(text string, priority int) {
	hash := crc32.ChecksumIEEE([]byte(text))
	for i := range *self {
		if (*self)[i].hash == hash {
			return // text is not new
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
