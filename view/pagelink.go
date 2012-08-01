package view

func NewPageLink(page **Page, title string) *PageLink {
	return &PageLink{Page: page, Title: title}
}

type PageLink struct {
	Page    **Page
	Content View   // If nil, then self.LinkTitle() will be used
	Title   string // If "", then self.Page.LinkTitle() will be used
	Rel     string
}

func (self *PageLink) URL(args ...string) string {
	return self.Page.URL(args...)
}

func (self *PageLink) LinkContent(urlArgs ...string) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(urlArgs...))
	}
	return self.Content
}

func (self *PageLink) LinkTitle(urlArgs ...string) string {
	if self.Title == "" {
		return self.Page.LinkTitle(urlArgs...)
	}
	return self.Title
}

func (self *PageLink) LinkRel() string {
	return self.Rel
}
