package view

import (
	"fmt"
)

type JQueryDialog struct {
	Div
	Title     string
	AutoOpen  bool
	CloseText string
}

func (self *JQueryDialog) Render(ctx *Context) (err error) {
	script := fmt.Sprintf(
		`jQuery(function(){
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
	return fmt.Sprintf(`jQuery("#%s").dialog("open");`, self.ID())
}

func (self *JQueryDialog) CloseScript() string {
	return fmt.Sprintf(`jQuery("#%s").dialog("close");`, self.ID())
}
