package view

// StringLink implements the LinkModel interface with string values
// for Url, Title, and Rel.
type StringLink struct {
	Url     string
	Content View   // If nil, then self.LinkTitle() will be used
	Title   string // If "", then self.URL will be used
	Rel     string
}

func (self *StringLink) URL(args ...string) string {
	return StringURL(self.Url).URL(args...)
}

func (self *StringLink) LinkContent(urlArgs ...string) View {
	if self.Content == nil {
		return HTML(self.LinkTitle())
	}
	return self.Content
}

func (self *StringLink) LinkTitle(urlArgs ...string) string {
	if self.Title == "" {
		return self.URL(urlArgs...)
	}
	return self.Title
}

func (self *StringLink) LinkRel() string {
	return self.Rel
}
