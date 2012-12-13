package view

import "fmt"

type ModalDialog struct {
	ViewBaseWithId
	Class   string
	Style   string
	Content View
}

func (self *ModalDialog) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *ModalDialog) Render(ctx *Context) (err error) {
	ctx.Response.RequireStyleURL("/css/avgrund.css", 0)
	ctx.Response.RequireScriptURL("/js/libs/avgrund.js", 0)
	ctx.Response.RequireScript(`jQuery(function(){jQuery("body").append("<div class='avgrund-cover'></div>");});`, 0)

	ctx.Response.XML.OpenTag("aside")
	ctx.Response.XML.Attrib("id", self.ID())
	ctx.Response.XML.Attrib("class", "avgrund-popup "+self.Class)
	ctx.Response.XML.AttribIfNotDefault("style", self.Style)
	if self.Content != nil {
		err = self.Content.Render(ctx)
	}
	ctx.Response.XML.CloseTagAlways()
	return err
}

func (self *ModalDialog) OpenScript() string {
	return fmt.Sprintf("Avgrund.show('#%s');", self.ID())
}

func (self *ModalDialog) OpenButton(text string) *Button {
	return &Button{Content: Escape(text), OnClick: self.OpenScript()}
}

const ModalDialogCloseScript = "Avgrund.hide();"

func ModalDialogCloseButton(text string) *Button {
	return &Button{Content: Escape(text), OnClick: ModalDialogCloseScript}
}
