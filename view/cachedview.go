package view

import (
	"time"
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
	data       []byte
	validUntil time.Time
	path       string
}

func (self *CachedView) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *CachedView) Render(response *Response) (err error) {
	if self.Content == nil {
		return nil
	}
	if true || Config.DisableCachedViews || len(response.Request.Params) > 0 || response.Request.Method != "GET" {
		return self.Content.Render(response)
	}
	if self.data == nil || time.Now().After(self.validUntil) {
		// todo cache headers
		// r := response.New()
		// err = self.Content.Render(r)
		// if err != nil {
		// 	return err
		// }
		// self.data = r.Bytes()
		// self.validUntil = time.Now().Add(self.Duration)
	}
	_, err = response.Write(self.data)
	return err
}

func (self *CachedView) ClearCache() {
	self.data = nil
}

func (self *CachedView) URL(response *Response) string {
	if viewWithURL, ok := self.Content.(ViewWithURL); ok {
		return viewWithURL.URL(response)
	}
	return StringURL(self.path).URL(response)
}

func (self *CachedView) SetPath(path string) {
	if viewWithURL, ok := self.Content.(ViewWithURL); ok {
		viewWithURL.SetPath(path)
	} else {
		self.path = path
	}
}
