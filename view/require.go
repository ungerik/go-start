package view

func RequireStyle(css string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireStyle(css, priority)
		return nil
	})
}

func RequireStyleURL(url string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireStyleURL(url, priority)
		return nil
	})
}

func RequireHeadScript(script string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireHeadScript(script, priority)
		return nil
	})
}

func RequireHeadScriptURL(url string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireHeadScriptURL(url, priority)
		return nil
	})
}

func RequireScript(script string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireScript(script, priority)
		return nil
	})
}

func RequireScriptURL(url string, priority int) View {
	return RenderView(func(ctx *Context) error {
		ctx.Response.RequireScriptURL(url, priority)
		return nil
	})
}
