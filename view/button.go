package view

///////////////////////////////////////////////////////////////////////////////
// Button

// Button represents a HTML input element of type button or submit.
type Button struct {
	ViewBaseWithId
	Name           string
	Value          interface{}
	Submit         bool
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
	response.XML.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	if self.Submit {
		response.XML.OpenTag("input").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
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
	} else {
		response.XML.OpenTag("button").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
		response.XML.Attrib("type", "button")
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
		if self.Content != nil {
			err = self.Content.Render(response)
		}
		response.XML.ForceCloseTag()
	}
	return nil
}

//func (self *Button) SetName(name string) {
//	self.Name = name
//	ViewChanged(self)
//}
//
//func (self *Button) SetValue(value string) {
//	self.Value = value
//	ViewChanged(self)
//}
