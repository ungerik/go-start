package view

///////////////////////////////////////////////////////////////////////////////
// Select

type Select struct {
	ViewBaseWithId
	Model    SelectModel
	Name     string
	Size     int // 0 shows all items, 1 shows a dropdownbox, other values show size items
	Class    string
	Disabled bool
}

func (self *Select) Render(ctx *Context) (err error) {
	ctx.Response.XML.OpenTag("select")
	ctx.Response.XML.AttribIfNotDefault("id", self.id)
	ctx.Response.XML.AttribIfNotDefault("class", self.Class)
	ctx.Response.XML.Attrib("name", self.Name)
	if self.Disabled {
		ctx.Response.XML.Attrib("disabled", "disabled")
	}

	if self.Model != nil {
		numOptions := self.Model.NumOptions()
		size := self.Size
		if size == 0 {
			size = numOptions
		}
		ctx.Response.XML.Attrib("size", size)

		for i := 0; i < numOptions; i++ {
			ctx.Response.XML.OpenTag("option")
			ctx.Response.XML.AttribIfNotDefault("value", self.Model.Value(i))
			if self.Model.Selected(i) {
				ctx.Response.XML.Attrib("selected", "selected")
			}
			if self.Model.Disabled(i) {
				ctx.Response.XML.Attrib("disabled", "disabled")
			}
			err = self.Model.RenderLabel(i, ctx)
			if err != nil {
				return err
			}
			ctx.Response.XML.CloseTag() // option
		}
	} else {
		ctx.Response.XML.Attrib("size", self.Size)
	}

	ctx.Response.XML.CloseTagAlways() // select
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Model

type SelectModel interface {
	NumOptions() int
	Value(index int) string
	Selected(index int) bool
	Disabled(index int) bool
	RenderLabel(index int, ctx *Context) (err error)
}

///////////////////////////////////////////////////////////////////////////////
// StringsSelectModel

type StringsSelectModel struct {
	Options        []string
	SelectedOption string
}

func (self *StringsSelectModel) NumOptions() int {
	return len(self.Options)
}

func (self *StringsSelectModel) Value(index int) string {
	return self.Options[index]
}

func (self *StringsSelectModel) Selected(index int) bool {
	return self.Options[index] == self.SelectedOption
}

func (self *StringsSelectModel) Disabled(index int) bool {
	return false
}

func (self *StringsSelectModel) RenderLabel(index int, ctx *Context) (err error) {
	ctx.Response.WriteString(self.Options[index])
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// IndexedStringsSelectModel

type IndexedStringsSelectModel struct {
	Options []string
	Index   int
}

func (self *IndexedStringsSelectModel) NumOptions() int {
	return len(self.Options)
}

func (self *IndexedStringsSelectModel) Value(index int) string {
	return self.Options[index]
}

func (self *IndexedStringsSelectModel) Selected(index int) bool {
	return index == self.Index
}

func (self *IndexedStringsSelectModel) Disabled(index int) bool {
	return false
}

func (self *IndexedStringsSelectModel) RenderLabel(index int, ctx *Context) (err error) {
	ctx.Response.WriteString(self.Options[index])
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// ValueLabelSelectModel

type ValueLabelSelectModel struct {
	ValuesAndLabels []string // Values and labels are interleaved starting with a value
	SelectedValue   string
}

func (self *ValueLabelSelectModel) NumOptions() int {
	return len(self.ValuesAndLabels) / 2
}

func (self *ValueLabelSelectModel) Value(index int) string {
	return self.ValuesAndLabels[index*2]
}

func (self *ValueLabelSelectModel) Selected(index int) bool {
	return self.ValuesAndLabels[index*2] == self.SelectedValue
}

func (self *ValueLabelSelectModel) Disabled(index int) bool {
	return false
}

func (self *ValueLabelSelectModel) RenderLabel(index int, ctx *Context) (err error) {
	ctx.Response.WriteString(self.ValuesAndLabels[index*2+1])
	return nil
}
