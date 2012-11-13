/*
The View of MVC.

Instead of using text templates, go-start creates HTML by using strongly typed Go composite literals.
For some HTML constructs there are high level Go representations,
but most HTML elements have a direct Go type or shortcut function counterpart.

Shortcut functions for HTML elements are all upper case,
other utility functions follow the normal camel case naming convention.

All shortcuts (see shortcuts.go):

	ViewOrError(view View, err error) View

	Escape(text string) HTML

	Printf(text string, args ...interface{}) HTML

	PrintfEscape(text string, args ...interface{}) HTML

	A(url interface{}, content ...interface{}) Link

	A_nofollow(url interface{}, content ...interface{}) Link

	A_blank(url interface{}, content ...interface{}) *Link

	A_blank_nofollow(url interface{}, content ...interface{}) *Link

	STYLE(css string) HTML

	StylesheetLink(url string) HTML

	SCRIPT(javascript string) HTML

	ScriptLink(url string) HTML

	RSSLink(title string, url URL) View

	IMG(url string, dimensions ...int)

	SECTION(class string, content ...interface{}) View

	DIV(class string, content ...interface{}) *Div

	SPAN(class string, content ...interface{}) *Span

	DivClearBoth() HTML

	CANVAS(class string, width, height int) *Canvas

	BR() HTML

	HR() HTML

	P(content ...interface{}) View

	H1(content ...interface{}) View

	H2(content ...interface{}) View

	H3(content ...interface{}) View

	H4(content ...interface{}) View

	H5(content ...interface{}) View

	H6(content ...interface{}) View

	B(content ...interface{}) View

	I(content ...interface{}) View

	Q(content ...interface{}) View

	DEL(content ...interface{}) View

	EM(content ...interface{}) View

	STRONG(content ...interface{}) View

	DFN(content ...interface{}) View

	CODE(content ...interface{}) View

	PRE(content ...interface{}) View

	ABBR(longTitle, abbreviation string) View

	UL(items ...interface{}) *List

	OL(items ...interface{}) *List

	NewView(content interface{}) View

	NewViews(contents ...interface{}) Views

	WrapContents(contents ...interface{}) View
*/
package view
