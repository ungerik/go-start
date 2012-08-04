package view

import "strings"

// StringURL implements the URL interface for a string.
type StringURL string

func (self StringURL) URL(response *Response) string {
	url := string(self)
	for _, arg := range response.Request.URLArgs {
		url = strings.Replace(url, PathFragmentPattern, arg, 1)
	}
	return url
}
