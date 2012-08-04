package view

// StringLink implements the LinkModel interface with string values
// for Url, Title, and Rel.
type StringLink struct {
	Url     string
	Content View   // If nil, then self.LinkTitle() will be used
	Title   string // If "", then self.URL will be used
	Rel     string
}

func (self *StringLink) URL(response *Response) string {
	return StringURL(self.Url).URL(response)
}

func (self *StringLink) LinkContent(response *Response) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(response))
	}
	return self.Content
}

func (self *StringLink) LinkTitle(response *Response) string {
	if self.Title == "" {
		return self.URL(response)
	}
	return self.Title
}

func (self *StringLink) LinkRel(response *Response) string {
	return self.Rel
}
