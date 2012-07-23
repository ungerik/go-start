package view

import "strings"

// StringURL implements the URL interface for a string.
type StringURL string

func (self StringURL) URL(args ...string) string {
	url := string(self)
	for _, arg := range args {
		url = strings.Replace(url, PathFragmentPattern, arg, 1)
	}
	return url
}
