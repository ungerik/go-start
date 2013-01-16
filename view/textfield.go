package view

type TextFieldType int

const (
	NormalTextField TextFieldType = iota
	PasswordTextField
	EmailTextField
	SearchTextField
)

///////////////////////////////////////////////////////////////////////////////
// TextField

type TextField struct {
	ViewBaseWithId
	Text        string
	Name        string
	Size        int
	MaxLength   int
	Type        TextFieldType
	TabIndex    int
	Class       string
	Placeholder string
	Title       string
	Readonly    bool
	Disabled    bool
	Required    bool   // HTML5
	Autofocus   bool   // HTML5
	Pattern     string // HTML5
}

func (self *TextField) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("input")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.AttribIfNotDefault("title", self.Title)

	ctx.Response.XML.Attrib("name", self.Name)
	ctx.Response.XML.AttribIfNotDefault("tabindex", self.TabIndex)
	ctx.Response.XML.AttribFlag("readonly", self.Readonly)
	ctx.Response.XML.AttribFlag("disabled", self.Disabled)
	ctx.Response.XML.AttribFlag("required", self.Required)
	ctx.Response.XML.AttribFlag("autofocus", self.Autofocus)
	ctx.Response.XML.AttribIfNotDefault("pattern", self.Pattern)

	switch self.Type {
	case PasswordTextField:
		ctx.Response.XML.Attrib("type", "password")
	case EmailTextField:
		ctx.Response.XML.Attrib("type", "email")
	case SearchTextField:
		ctx.Response.XML.Attrib("type", "search")
	default:
		ctx.Response.XML.Attrib("type", "text")
	}

	ctx.Response.XML.AttribIfNotDefault("size", self.Size)
	ctx.Response.XML.AttribIfNotDefault("maxlength", self.MaxLength)
	ctx.Response.XML.AttribIfNotDefault("placeholder", self.Placeholder)

	ctx.Response.XML.Attrib("value", self.Text)

	ctx.Response.XML.CloseTag()
	return nil
}

func (self *TextField) SetRequired(required bool) {
	self.Required = required
}
