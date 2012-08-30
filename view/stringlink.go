package view

// StringLink implements the LinkModel interface with string values
// for Url, Title, and Rel.
type StringLink struct {
	Url     string
	Content View   // If nil, then self.LinkTitle() will be used
	Title   string // If "", then self.URL will be used
	Rel     string
}

func (self *StringLink) URL(ctx *Context) string {
	return StringURL(self.Url).URL(ctx)
}

func (self *StringLink) LinkContent(ctx *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(ctx))
	}
	return self.Content
}

func (self *StringLink) LinkTitle(ctx *Context) string {
	if self.Title == "" {
		return self.URL(ctx)
	}
	return self.Title
}

func (self *StringLink) LinkRel(ctx *Context) string {
	return self.Rel
}
