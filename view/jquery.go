package view

import (
	"fmt"
)

var JQuery HTML = `<script src="//ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js"></script><script>window.jQuery || document.write('<script src="/js/libs/jquery-1.7.1.min.js"><\/script>')</script>`
var JQueryUI HTML = `<script src="//ajax.googleapis.com/ajax/libs/jqueryui/1.8.17/jquery-ui.min.js"></script><script>window.jQuery.ui || document.write('<script src="/js/libs/jquery-ui-1.8.17.custom.min.js"><\/script>')</script>`

func RequireJQuery(priority int) View {
	return RequireScript(string(JQuery), priority)
}

func RequireJQueryUI(priority int) View {
	return RequireScript(string(JQueryUI), priority)
}

/*

Add the following to your CSS to set the max height of the drop-down:

	.ui-autocomplete {
	    max-height: 200px;
	    overflow-y: auto;
	    overflow-x: hidden;
	}
*/

func JQueryUIAutocompleteFromURL(domSelector string, dataURL URL, minLength int) View {
	return RenderView(
		func(ctx *Context) (err error) {
			ctx.Response.Printf("<script>jQuery('%s').autocomplete({source:'%s',minLength:%d});</script>", domSelector, dataURL.URL(ctx), minLength)
			return nil
		},
	)
}

func RequireJQueryUIAutocompleteFromURL(domSelector string, dataURL URL, minLength int, priority int) View {
	return RenderView(
		func(ctx *Context) (err error) {
			script := fmt.Sprintf("jQuery('%s').autocomplete({source:'%s',minLength:%d});", domSelector, dataURL.URL(ctx), minLength)
			ctx.Response.RequireScript(script, priority)
			return nil
		},
	)
}

func JQueryUIAutocomplete(domSelector string, options []string, minLength int) View {
	return RenderView(
		func(ctx *Context) (err error) {
			ctx.Response.Printf("<script>jQuery('%s').autocomplete({source:[", domSelector)
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

// Dependencies of jQuery UI Dialog are not included.
type JQueryDialog struct {
	Div
	Title     string
	AutoOpen  bool
	CloseText string
}

func (self *JQueryDialog) Render(ctx *Context) (err error) {
	script := fmt.Sprintf(
		`jQuery(function() {
			jQuery("#%s").dialog({
				title: %s,
				autoOpen: %t,
				closeText: %s,
				modal: true
			});
		})`,
		self.ID(),
		self.Title,
		self.AutoOpen,
		self.CloseText,
	)
	ctx.Response.RequireScript(script, 0)
	return self.Div.Render(ctx)
}

func (self *JQueryDialog) OpenScript() string {
	return fmt.Sprintf("jQuery('#%s').dialog('open');", self.ID())
}

func (self *JQueryDialog) CloseScript() string {
	return fmt.Sprintf("jQuery('#%s').dialog('close');", self.ID())
}
