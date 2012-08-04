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

func (self *PageLink) URL(response *Response) string {
	return self.Page.URL(response)
}

func (self *PageLink) LinkContent(response *Response) View {
	if self.Content == nil {
		return HTML(self.LinkTitle(response))
	}
	return self.Content
}

func (self *PageLink) LinkTitle(response *Response) string {
	if self.Title == "" {
		return self.URL(response)
	}
	return self.Title
}

func (self *PageLink) LinkRel(response *Response) string {
	return self.Rel
}
