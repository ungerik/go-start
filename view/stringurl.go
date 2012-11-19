package view

import "strings"

// StringURL implements the URL interface for a URL string.
type StringURL string

func (self StringURL) URL(ctx *Context) string {
	url := string(self)
	for _, arg := range ctx.URLArgs {
		url = strings.Replace(url, PathFragmentPattern, arg, 1)
	}
	return ctx.Request.AddProtocolAndHostToURL(url)
}
