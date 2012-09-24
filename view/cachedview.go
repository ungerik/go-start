package view

import (
	"net/http"
	"time"

	// "github.com/ungerik/go-start/debug"
)

var cachedViews map[string]*CachedView

// CacheView caches view for duration by wrapping it with a CachedView object.
// The CachedView is added to an internal list so that ClearAllCaches()
// can call ClearCache() on every CachedView.
// UncacheView() removes the CachedView from the internal list.
func CacheView(duration time.Duration, view View) *CachedView {
	cachedView := &CachedView{
		Content:  view,
		Duration: duration,
	}
	cachedView.Init(cachedView) // Acquire an id
	if cachedViews == nil {
		cachedViews = make(map[string]*CachedView)
	}
	cachedViews[cachedView.id] = cachedView
	return cachedView
}

// UncacheView removes the CachedView from the internal list CachedViews
// that is used by ClearAllCaches().
// It returnes the View wrapped CachedView.
func UncacheView(cachedView *CachedView) View {
	delete(cachedViews, cachedView.id)
	return cachedView.Content
}

// ClearAllCaches cals ClearCache() for all CachedView objects
// created with  CacheView().
func ClearAllCaches() {
	for _, cachedView := range cachedViews {
		cachedView.ClearCache()
	}
}

// CachedView implements ViewWithURL. If Content implements ViewWithURL too,
// the ViewWithURL methods will be forwarded to Content.
// Else CachedView provides its own implementation of ViewWithURL.
// Caching will be disabled, when the request is not a GET or has url parameters.
type CachedView struct {
	ViewBaseWithId
	Content    View
	Duration   time.Duration
	path       string
	header     http.Header
	body       []byte
	validUntil time.Time
}

func (self *CachedView) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *CachedView) Render(ctx *Context) (err error) {
	if self.Content == nil {
		return nil
	}
	if Config.DisableCachedViews || len(ctx.Request.Params) > 0 || ctx.Request.Method != "GET" {
		return self.Content.Render(ctx)
	}
	if self.body == nil || time.Now().After(self.validUntil) {
		err = self.Content.Render(ctx)
		if err != nil {
			return err
		}
		self.header = ctx.Response.Header()
		self.body = ctx.Response.Bytes()
		self.validUntil = time.Now().Add(self.Duration)
		return nil
	}

	// debug.Print("Responding with cached version of " + ctx.Request.URLString())
	for key, value := range self.header {
		ctx.Response.Header()[key] = value
	}
	_, err = ctx.Response.Write(self.body)
	return err
}

func (self *CachedView) ClearCache() {
	self.header = nil
	self.body = nil
}

func (self *CachedView) URL(ctx *Context) string {
	if viewWithURL, ok := self.Content.(ViewWithURL); ok {
		return viewWithURL.URL(ctx)
	}
	return StringURL(self.path).URL(ctx)
}

func (self *CachedView) SetPath(path string) {
	self.path = path
	if viewWithURL, ok := self.Content.(ViewWithURL); ok {
		viewWithURL.SetPath(path)
	}
}
