package view

// URLLink implements the LinkModel interface with a reference to a URL
// and string values for Title and Rel.
type URLLink struct {
	Url     URL
	Content View   // If is nil, then self.LinkTitle() will be used
	Title   string // If is "", then self.URL will be used
	Rel     string
}

func (self *URLLink) URL(args ...string) string {
	return self.Url.URL(args...)
}

func (self *URLLink) LinkContent(response *Response) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(response))
	}
	return self.Content
}

func (self *URLLink) LinkTitle(response *Response) string {
	if self.Title == "" {
		return self.Url.URL(response.Request.URLArgs...)
	}
	return self.Title
}

func (self *URLLink) LinkRel(response *Response) string {
	return self.Rel
}
