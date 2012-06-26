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

func (self *StringLink) LinkContent(context *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(context))
	}
	return self.Content
}

func (self *StringLink) LinkTitle(context *Context) string {
	if self.Title == "" {
		return self.URL(context.PathArgs...)
	}
	return self.Title
}

func (self *StringLink) LinkRel(context *Context) string {
	return self.Rel
}
