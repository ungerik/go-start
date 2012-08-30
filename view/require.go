package view

// RequireStyle adds dynamic CSS content to the page.
// Multiple dynamic entries will be sorted by priority.
// Dynamic CSS will be inserted after the regular CSS of the page.
// 
// Use this feature to dynamically add CSS to the page if the
// HTML content requires it.
func RequireStyle(css string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireStyle(css, priority)
		return nil
	})
}

// RequireStyleURL adds a dynamic CSS link to the page.
// Multiple dynamic entries will be sorted by priority.
// Dynamic CSS will be inserted after the regular CSS of the page.
// 
// Use this feature to dynamically add CSS to the page if the
// HTML content requires it.
func RequireStyleURL(url string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireStyleURL(url, priority)
		return nil
	})
}

// RequireHeadScript adds dynamic JavaScript to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// head-scripts of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func RequireHeadScript(script string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireHeadScript(script, priority)
		return nil
	})
}

// RequireHeadScriptURL adds a dynamic JavaScript link to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// head-scripts of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func RequireHeadScriptURL(url string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireHeadScriptURL(url, priority)
		return nil
	})
}

// RequireScript adds dynamic JavaScript to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// scripts near the end of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func RequireScript(script string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireScript(script, priority)
		return nil
	})
}

// RequireScriptURL adds a dynamic JavaScript link to the page.
// Multiple dynamic entries will be sorted by priority.
// The dynamic JavaScript will be inserted after the regular
// scripts near the end of the page.
// 
// Use this feature to dynamically add JavaScript to
// the page if the HTML content requires it.
func RequireScriptURL(url string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireScriptURL(url, priority)
		return nil
	})
}
