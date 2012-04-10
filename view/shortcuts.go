package view

import (
	"fmt"
	"html"
	"reflect"
)

///////////////////////////////////////////////////////////////////////////////
// Shortcuts

// Escape HTML escapes a string.
func Escape(text string) HTML {
	return HTML(html.EscapeString(text))
}

// Printf creates an unescaped HTML string.
func Printf(text string, args ...interface{}) HTML {
	return HTML(fmt.Sprintf(text, args...))
}

// PrintfEscape creates an escaped HTML string.
func PrintfEscape(text string, args ...interface{}) HTML {
	return Escape(fmt.Sprintf(text, args...))
}

// A creates a HTML link.
func A(url interface{}, content ...interface{}) *Link {
	return &Link{Model: NewLinkModel(url, content...)}
}

// A_blank creates a HTML link with target="_blank"
func A_blank(url interface{}, content ...interface{}) *Link {
	return &Link{NewWindow: true, Model: NewLinkModel(url, content...)}
}

// Img creates a HTML img element for an URL with optional width and height.
// The first int of dimensions is width, the second one height.
func IMG(url string, dimensions ...int) View {
	var width int
	var height int
	dimCount := len(dimensions)
	if dimCount >= 1 {
		width = dimensions[0]
		if dimCount >= 2 {
			height = dimensions[1]
		}
	}
	if url == "" && width > 0 && height > 0 {
		return &DummyImage{Width: width, Height: height}
	}
	return &Image{URL: url, Width: width, Height: height}
}

func SECTION(class string, content ...interface{}) View {
	return &ShortTag{Tag: "section", Class: class, Content: WrapContents(content...)}
}

// Creates a Div object with a HTML class attribute and optional content.
func DIV(class string, content ...interface{}) *Div {
	return &Div{Class: class, Content: WrapContents(content...)}
}

func DivClearBoth() HTML {
	return HTML("<div style='clear:both'></div>")
}

func CANVAS(class string, width, height int) *Canvas {
	return &Canvas{Class: class, Width: width, Height: height}
}

func BR() HTML {
	return HTML("<br/>")
}

func HR() HTML {
	return HTML("<hr/>")
}

func P(content ...interface{}) View {
	return &ShortTag{Tag: "p", Content: WrapContents(content...)}
}

func H1(content ...interface{}) View {
	return &ShortTag{Tag: "h1", Content: WrapContents(content...)}
}

func H2(content ...interface{}) View {
	return &ShortTag{Tag: "h2", Content: WrapContents(content...)}
}

func H3(content ...interface{}) View {
	return &ShortTag{Tag: "h3", Content: WrapContents(content...)}
}

func H4(content ...interface{}) View {
	return &ShortTag{Tag: "h4", Content: WrapContents(content...)}
}

func H5(content ...interface{}) View {
	return &ShortTag{Tag: "h5", Content: WrapContents(content...)}
}

func H6(content ...interface{}) View {
	return &ShortTag{Tag: "h6", Content: WrapContents(content...)}
}

func B(content ...interface{}) View {
	return &ShortTag{Tag: "b", Content: WrapContents(content...)}
}

func I(content ...interface{}) View {
	return &ShortTag{Tag: "i", Content: WrapContents(content...)}
}

func Q(content ...interface{}) View {
	return &ShortTag{Tag: "q", Content: WrapContents(content...)}
}

func DEL(content ...interface{}) View {
	return &ShortTag{Tag: "del", Content: WrapContents(content...)}
}

func EM(content ...interface{}) View {
	return &ShortTag{Tag: "em", Content: WrapContents(content...)}
}

func STRONG(content ...interface{}) View {
	return &ShortTag{Tag: "strong", Content: WrapContents(content...)}
}

func DFN(content ...interface{}) View {
	return &ShortTag{Tag: "dfn", Content: WrapContents(content...)}
}

func CODE(content ...interface{}) View {
	return &ShortTag{Tag: "code", Content: WrapContents(content...)}
}

func PRE(content ...interface{}) View {
	return &ShortTag{Tag: "pre", Content: WrapContents(content...)}
}

func ABBR(longTitle, abbreviation string) View {
	return &ShortTag{Tag: "abbr", Attribs: map[string]string{"title": longTitle}, Content: Escape(abbreviation)}
}

// Ul is a shortcut to create an unordered list by wrapping items as HTML views.
// NewView will be called for every passed item.
func UL(items ...interface{}) *List {
	model := make(ViewsListModel, len(items))
	for i, item := range items {
		model[i] = NewView(item)
	}
	return &List{Model: model}
}

// Ul is a shortcut to create an ordered list by wrapping items as HTML views.
// NewView will be called for every passed item.
func OL(items ...interface{}) *List {
	list := UL(items...)
	list.Ordered = true
	return list
}

// Encapsulates content as View.
// Strings or fmt.Stringer implementations will be HTML escaped.
// View implementations will be passed through. 
func NewView(content interface{}) View {
	if content == nil {
		return nil
	}
	if view, ok := content.(View); ok {
		return view
	}
	if stringer, ok := content.(fmt.Stringer); ok {
		return Escape(stringer.String())
	}
	v := reflect.ValueOf(content)
	if v.Kind() != reflect.String {
		panic(fmt.Errorf("Invalid content type: %T (must be gostart/view.View, fmt.Stringer or a string)", content))
	}
	return Escape(v.String())
}

// Encapsulates multiple content arguments as views by calling NewView() for every one of them.
func NewViews(contents ...interface{}) Views {
	count := len(contents)
	if count == 0 {
		return nil
	}
	views := make(Views, count)
	for i, content := range contents {
		views[i] = NewView(content)
	}
	return views
}

// Encapsulates multiple content arguments as View by calling NewView() for every one of them.
// It is more efficient for one view because the view is passed through instead of wrapped
// with a Views slice like NewViews does.
func WrapContents(contents ...interface{}) View {
	count := len(contents)
	switch count {
	case 0:
		return nil
	case 1:
		return NewView(contents[0])
	}
	views := make(Views, count)
	for i, content := range contents {
		views[i] = NewView(content)
	}
	return views
}
