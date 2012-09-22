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

func newResponse(webContext *web.Context) *Response {
	response := &Response{
		webContext: webContext,
	}
	response.PushBody()
	return response
}

type responseBody struct {
	buffer *bytes.Buffer
	xml    *utils.XMLWriter
}

type Response struct {
	webContext *web.Context

	Session *Session

	bodyStack []responseBody
	// XML allowes the Response to be used as utils.XMLWriter
	XML *utils.XMLWriter

	dynamicStyle       dependencyHeap
	dynamicHeadScripts dependencyHeap
	dynamicScripts     dependencyHeap
}

// PushBody pushes the buffer of the response body on a stack
// and sets a new empty buffer.
// This can be used to render intermediate text results.
// Note: Only the response body is pushed, all other state changes
// like setting headers will affect the final response.
func (self *Response) PushBody() {
	var b responseBody
	b.buffer = new(bytes.Buffer)
	b.xml = utils.NewXMLWriter(b.buffer)
	self.bodyStack = append(self.bodyStack, b)
	self.XML = b.xml
}

// PopBody pops the buffer of the response body from the stack
// and returns its content.
func (self *Response) PopBody() (bufferData []byte) {
	last := len(self.bodyStack) - 1
	bufferData = self.bodyStack[last].buffer.Bytes()
	self.bodyStack = self.bodyStack[0:last]
	self.XML = self.bodyStack[last-1].xml
	return bufferData
}

// PopBodyString pops the buffer of the response body from the stack
// and returns its content as string.
func (self *Response) PopBodyString() (bufferData string) {
	return string(self.PopBody())
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
	return self.bodyStack[len(self.bodyStack)-1].buffer.String()
}

func (self *Response) Bytes() []byte {
	return self.bodyStack[len(self.bodyStack)-1].buffer.Bytes()
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

func (self *Response) SetContentTypeByExt(ext string) {
	self.webContext.ContentType(ext)
}

// todo gets overwritten bei web.go, so use SetContentTypeByExt
// func (self *Response) SetContentType(ctype string) {
// 	self.Header().Add("Content-Type", ctype)
// }

// ContentDispositionAttachment makes the webbrowser open a
// "Save As.." dialog for the response.
func (self *Response) ContentDispositionAttachment(filename string) {
	self.Header().Add("Content-Type", "application/x-unknown")
	self.Header().Add("Content-Disposition", "attachment;filename="+filename)
}

// RequireStyle adds dynamic CSS content to the page.
// Multiple dynamic entries will be sorted by priority.
// Dynamic CSS will be inserted after the regular CSS of the page.
// 
// Use this feature to dynamically add CSS to the page if the
// HTML content requires it.
func (self *Response) RequireStyle(css string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = newDependencyHeap()
	}
	self.dynamicStyle.AddIfNew("<style>"+css+"</style>", priority)
}

// RequireStyleURL adds a dynamic CSS link to the page.
// Multiple dynamic entries will be sorted by priority.
// Dynamic CSS will be inserted after the regular CSS of the page.
// 
// Use this feature to dynamically add CSS to the page if the
// HTML content requires it.
func (self *Response) RequireStyleURL(url string, priority int) {
	if self.dynamicStyle == nil {
		self.dynamicStyle = newDependencyHeap()
	}
	self.dynamicStyle.AddIfNew("<link rel='stylesheet' href='"+url+"'>", priority)
}

// RequireHeadScript adds dynamic JavaScript to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// head-scripts of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) RequireHeadScript(script string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = newDependencyHeap()
	}
	self.dynamicHeadScripts.AddIfNew("<script>"+script+"</script>", priority)
}

// RequireHeadScriptURL adds a dynamic JavaScript link to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// head-scripts of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) RequireHeadScriptURL(url string, priority int) {
	if self.dynamicHeadScripts == nil {
		self.dynamicHeadScripts = newDependencyHeap()
	}
	self.dynamicHeadScripts.AddIfNew("<script src='"+url+"'></script>", priority)
}

// RequireScript adds dynamic JavaScript to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// scripts near the end of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) RequireScript(script string, priority int) {
	if self.dynamicScripts == nil {
		self.dynamicScripts = newDependencyHeap()
	}
	self.dynamicScripts.AddIfNew("<script>"+script+"</script>", priority)
}

// RequireScriptURL adds a dynamic JavaScript link to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// scripts near the end of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func (self *Response) RequireScriptURL(url string, priority int) {
	if self.dynamicScripts == nil {
		self.dynamicScripts = newDependencyHeap()
	}
	self.dynamicScripts.AddIfNew("<script src='"+url+"'></script>", priority)
}

///////////////////////////////////////////////////////////////////////////////
// dependencyHeap

func newDependencyHeap() dependencyHeap {
	dh := make(dependencyHeap, 0, 1)
	heap.Init(&dh)
	return dh
}

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
