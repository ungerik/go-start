package view

func NewURLWithArgs(url URL, args ...string) *URLWithArgs {
	return &URLWithArgs{
		Url:  url,
		Args: args,
	}
}

// URLWithArgs binds Args to an URL.
// Url.URL() will be called with response.URLArgs(Args...)
type URLWithArgs struct {
	Url  URL
	Args []string
}

func (self *URLWithArgs) URL(ctx *Context) string {
	return self.Url.URL(ctx.ForURLArgs(self.Args...))
}
