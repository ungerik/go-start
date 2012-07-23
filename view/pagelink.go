package view

func NewPageLink(page **Page, title string) *PageLink {
	return &PageLink{Page: page, Title: title}
}

type PageLink struct {
	Page    **Page
	Content View   // If nil, then self.LinkTitle() will be used
	Title   string // If "", then self.Page.LinkTitle(context) will be used
	Rel     string
}

func (self *PageLink) URL(args ...string) string {
	return self.Page.URL(args...)
}

func (self *PageLink) LinkContent(context *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(context))
	}
	return self.Content
}

func (self *PageLink) LinkTitle(context *Context) string {
	if self.Title == "" {
		return self.Page.LinkTitle(context)
	}
	return self.Title
}

func (self *PageLink) LinkRel(context *Context) string {
	return self.Rel
}
