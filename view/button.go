package view

type Button struct {
	ViewBaseWithId
	Name           string
	Class          string
	Disabled       bool
	TabIndex       int
	OnClick        string
	OnClickConfirm string // Will add a confirmation dialog for onclick
	Content        View   // Only used when Submit is false
}

func (self *Button) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Button) Render(response *Response) (err error) {
	response.XML.OpenTag("button")
	response.XML.AttribIfNotDefault("id", self.id)
	response.XML.AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("type", "button")
	response.XML.AttribIfNotDefault("name", self.Name)
	if self.Disabled {
		response.XML.Attrib("disabled", "disabled")
	}
	response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.OnClickConfirm != "" {
		response.XML.Attrib("onclick", "return confirm('", self.OnClickConfirm, "');")
	} else {
		response.XML.AttribIfNotDefault("onclick", self.OnClick)
	}
	if self.Content != nil {
		err = self.Content.Render(response)
	}
	response.XML.ForceCloseTag()
	return err
}
