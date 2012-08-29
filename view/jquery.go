package view

var (
	JQuery   HTML = `<script src="//ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js"></script><script>window.jQuery || document.write('<script src="/js/libs/jquery-1.7.1.min.js"><\/script>')</script>`
	JQueryUI HTML = `<script src="//ajax.googleapis.com/ajax/libs/jqueryui/1.8.17/jquery-ui.min.js"></script><script>window.jQuery.ui || document.write('<script src="/js/libs/jquery-ui-1.8.17.custom.min.js"><\/script>')</script>`
)

func JQueryUIAutocompleteFromURL(domSelector string, dataURL URL, minLength int) View {
	return RenderView(
		func(ctx *Context) (err error) {
			url := dataURL.URL(ctx)
			ctx.Response.Printf("<script>$('%s').autocomplete({source:'%s',minLength:%d});</script>", domSelector, url, minLength)
			return nil
		},
	)
}

func JQueryUIAutocomplete(domSelector string, options []string, minLength int) View {
	return RenderView(
		func(ctx *Context) (err error) {
			ctx.Response.Printf("<script>$('%s').autocomplete({source:[", domSelector)
			for i := range options {
				if i > 0 {
					ctx.Response.WriteByte(',')
				}
				ctx.Response.WriteByte('"')
				ctx.Response.WriteString(options[i])
				ctx.Response.WriteByte('"')
			}
			ctx.Response.Printf("],minLength:%d});</script>", minLength)
			return nil
		},
	)
}
