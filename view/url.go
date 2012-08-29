package view

// URL is an interface to return URL strings depending on the request path args.
type URL interface {
	URL(ctx *Context) string
}
