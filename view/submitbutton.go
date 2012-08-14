package view

type SubmitButton struct {
	ViewBaseWithId
	Name           string
	Value          interface{}
	Class          string
	Disabled       bool
	TabIndex       int
	OnClick        string
	OnClickConfirm string // Will add a confirmation dialog for onclick
}

func (self *SubmitButton) Render(response *Response) (err error) {
	response.XML.OpenTag("input")
	response.XML.AttribIfNotDefault("id", self.id)
	response.XML.AttribIfNotDefault("class", self.Class)
	response.XML.Attrib("type", "submit")
	response.XML.AttribIfNotDefault("name", self.Name)
	response.XML.AttribIfNotDefault("value", self.Value)
	if self.Disabled {
		response.XML.Attrib("disabled", "disabled")
	}
	response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	if self.OnClickConfirm != "" {
		response.XML.Attrib("onclick", "return confirm('", self.OnClickConfirm, "');")
	} else {
		response.XML.AttribIfNotDefault("onclick", self.OnClick)
	}
	response.XML.CloseTag()
	return nil
}
