package view

func NewPageLink(page **Page, title string) *PageLink {
	return &PageLink{Page: page, Title: title}
}

type PageLink struct {
	Page    **Page
	Content View
	Title   string
	Rel     string
}

func (self *PageLink) URL(ctx *Context) string {
	return self.Page.URL(ctx)
}

func (self *PageLink) LinkContent(ctx *Context) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(ctx))
	}
	return self.Content
}

func (self *PageLink) LinkTitle(ctx *Context) string {
	if self.Title == "" {
		return self.Page.LinkTitle(ctx)
	}
	return self.Title
}

func (self *PageLink) LinkRel(ctx *Context) string {
	return self.Rel
}
