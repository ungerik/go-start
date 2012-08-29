package view

// URLLink implements the LinkModel interface with a reference to a URL
// and string values for Title and Rel.
type URLLink struct {
	Url     URL
	Content View   // If is nil, then self.LinkTitle() will be used
	Title   string // If is "", then self.URL will be used
	Rel     string
}

func (self *URLLink) URL(ctx *Context) string {
	return self.Url.URL(ctx)
}

func (self *URLLink) LinkContent(ctx *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(ctx))
	}
	return self.Content
}

func (self *URLLink) LinkTitle(ctx *Context) string {
	if self.Title == "" {
		return self.URL(ctx)
	}
	return self.Title
}

func (self *URLLink) LinkRel(ctx *Context) string {
	return self.Rel
}
