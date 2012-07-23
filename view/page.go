package view

import (
	"github.com/ungerik/go-start/templatesystem"
	// "github.com/ungerik/go-start/utils"
	"html"
)

<<<<<<< HEAD
=======
// PageWriteFunc is used by Page to write dynamic or static content to the page.
type PageWriteFunc func(context *Context, writer io.Writer) (err error)

// PageTitle writes a static page title.
func PageTitle(title string) PageWriteFunc {
	return PageWrite(title)
}

// PageMetaDescription writes a static meta description.
func PageMetaDescription(description string) PageWriteFunc {
	return PageWrite(description)
}

// PageWrite writes static text.
func PageWrite(text string) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		writer.Write([]byte(text))
		return nil
	}
}

// PageWriters combines multiple PageWriteFunc into a single PageWriteFunc
// by calling them one after another.
func PageWriters(funcs ...PageWriteFunc) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		for _, f := range funcs {
			if f != nil {
				if err = f(context, writer); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// IndirectPageWriter takes the pointer to a PageWriteFunc variable
// and dereferences it when the returned PageWriteFunc is called.
// Used to break dependency cycles of variable initializations by
// using a pointer to a variable instead of its value.
func IndirectPageWriter(pageWritePtr *PageWriteFunc) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		return (*pageWritePtr)(context, writer)
	}
}

// PageWritersFilterPort calls funcs only
// if the request is made to a specific port
func PageWritersFilterPort(port uint16, funcs ...PageWriteFunc) PageWriteFunc {
	if len(funcs) == 0 {
		return nil
	}
	return func(context *Context, writer io.Writer) (err error) {
		if context.RequestPort() != port {
			return nil
		}
		return PageWriters(funcs...)(context, writer)
	}
}

// Stylesheet writes a HTML style tag with the passed css as content.
func Stylesheet(css string) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		writer.Write([]byte("<style>"))
		writer.Write([]byte(css))
		writer.Write([]byte("</style>\n"))
		return nil
	}
}

// StylesheetURL writes a HTML style tag with that references url.
func StylesheetURL(url string) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		writer.Write([]byte("<link rel='stylesheet' href='"))
		writer.Write([]byte(url))
		writer.Write([]byte("'>\n"))
		return nil
	}
}

// Script writes a HTML script tag with the passed script as content.
func Script(script string) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		writer.Write([]byte("<script>"))
		writer.Write([]byte(script))
		writer.Write([]byte("</script>\n"))
		return nil
	}
}

// ScriptURL writes a HTML script tag with that references url.
func ScriptURL(url string) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		writer.Write([]byte("<script src='"))
		writer.Write([]byte(url))
		writer.Write([]byte("'></script>\n"))
		return nil
	}
}

// RSS a application/rss+xml link tag with the given title and url.
func RSS(title string, url URL) PageWriteFunc {
	return func(context *Context, writer io.Writer) (err error) {
		writer.Write([]byte("<link rel='alternate' type='application/rss+xml' title='"))
		writer.Write([]byte(title))
		writer.Write([]byte("' href='"))
		writer.Write([]byte(url.URL(context.PathArgs...)))
		writer.Write([]byte("'>\n"))
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// Page

>>>>>>> master
/*
Page is the basis to render complete HTML pages.
An arbitrary View or ViewWithURL can be used to render other text formats
like CSS or JSON.

A HTML5 Boilerplate template is used by default. See:
	* gostart/templates/html5boilerplate.html
	* gostart/static/css/html5boilerplate/normalize.css
	* gostart/static/css/html5boilerplate/poststyle.css

Note: In the current version Mustache is always used as templates system.
This will be changed to the Go v1 template system in the Go v1 syntax release.

Most HTML head and script specific text is written by Renderer functions
that receive the request context as an argument.
That way the content can be created dynamically.

Wrapper functions for static content are provided for convenience.
See functions under Renderer below.

Example:

	&Page{WriteTitle: func(response *Response, writer io.Writer) (err error) {
		writer.Write([]byte("Could be a dynamic title"))
		return nil
	}}

	&Page{WriteTitle: PageTitle("Static Title")}


To avoid reading the same data multiple times from the database in Renderer
or dynamic views in the content structure, OnPreRender can be used to
query and set page wide data only once at the request context.Data.

Example:

	&Page{
		OnPreRender: func(page *Page, response *Response) (err error) {
			context.Data = &MyPerPageData{SomeText: "Hello World!"}
		},
		Content: DynamicView(
			func(response *Response) (view View, err error) {
				myPerPageData := context.Data.(*MyPerPageData)
				return HTML(myPerPageData.SomeText), nil
			},
		),
	}

*/
type Page struct {
	Template

	// Called before any other function when rendering the page
	OnPreRender func(page *Page, response *Response) (err error)

	// Writes the head title tag
	Title Renderer

	// Writes the head meta description tag
	MetaDescription Renderer

	// Content of the head meta viewport tag,
	// Config.Page.DefaultMetaViewport will be used if ""
	MetaViewport string

	// Write additional HTML head content
	AdditionalHead Renderer

	// Write head content before the stylesheet link
	PreCSS Renderer

	// stylesheet link URL
	CSS URL

	// Write head content after the stylesheet link
	PostCSS Renderer

	// Write scripts as last element of the HTML head
	HeadScripts Renderer

	// HTML body content. Will be wrapped by a div with class="container"
	Content View

	// Write scripts after body content
	Scripts Renderer

	path string

	// That way of linking to favicons my be removed in the future:
	Favicon16x16URL   string
	Favicon57x57URL   string
	Favicon72x72URL   string
	Favicon114x114URL string
	Favicon129x129URL string
}

func (self *Page) Init(thisView View) {
	if self == nil {
		panic("Page is nil")
	}
	if thisView == self.thisView {
		return // already initialized
	}

	// Uses alwaysMustache template rendering
	self.TemplateSystem = &templatesystem.Mustache{}
	if self.Filename == "" {
		self.Filename = Config.Page.Template
	}
	self.Template.Init(thisView)
}

func (self *Page) IterateChildren(callback IterateChildrenCallback) {
	if self.Content != nil {
		callback(self, self.Content)
	}
}

func (self *Page) SetPath(path string) {
	self.path = path
}

// Implements the URL and LinkModel interface
func (self *Page) URL(args ...string) string {
	return StringURL(self.path).URL(args...)
}

// Implements the LinkModel interface
func (self *Page) LinkContent(response *Response) View {
	return HTML(self.LinkTitle(response))
}

// Implements the LinkModel interface
func (self *Page) LinkTitle(response *Response) string {
	if self.Title == nil {
		return ""
	}
	r := response.New()
	err := self.Title.Render(r)
	if err != nil {
		//return err.String()
		panic(err)
	}
	return r.String()
}

// Implements the LinkModel interface
func (self *Page) LinkRel(response *Response) string {
	return ""
}

func (self *Page) Render(response *Response) (err error) {
	if self.OnPreRender != nil {
		err = self.OnPreRender(self, response)
		if err != nil {
			return err
		}
	}

	var templateContext struct {
		Title              string
		MetaDescription    string
		MetaViewport       string
		Head               string
		PreCSS             string
		CSS                string
		PostCSS            string
		DynamicStyle       string
		HeadScripts        string
		DynamicHeadScripts string
		Scripts            string
		DynamicScripts     string
		Favicon16x16URL    string
		Favicon57x57URL    string
		Favicon72x72URL    string
		Favicon114x114URL  string
		Favicon129x129URL  string
		Content            string
	}

	if self.Title != nil {
		r := response.New()
		err := self.Title.Render(r)
		if err != nil {
			return err
		}
		templateContext.Title = html.EscapeString(r.String())
	}
	if self.MetaDescription != nil {
		r := response.New()
		err := self.MetaDescription.Render(r)
		if err != nil {
			return err
		}
		templateContext.MetaDescription = html.EscapeString(r.String())
	}

	metaViewport := self.MetaViewport
	if metaViewport == "" {
		metaViewport = Config.Page.DefaultMetaViewport
	}
	templateContext.MetaViewport = metaViewport

	additionalHead := self.AdditionalHead
	if additionalHead == nil {
		additionalHead = Config.Page.DefaultAdditionalHead
	}
	if additionalHead != nil {
		r := response.New()
		err := additionalHead.Render(r)
		if err != nil {
			return err
		}
		templateContext.Head = r.String()
	}

	//templateContext.Meta = self.Meta
	templateContext.Favicon16x16URL = self.Favicon16x16URL
	templateContext.Favicon57x57URL = self.Favicon57x57URL
	templateContext.Favicon72x72URL = self.Favicon72x72URL
	templateContext.Favicon114x114URL = self.Favicon114x114URL
	templateContext.Favicon129x129URL = self.Favicon129x129URL

	if self.PreCSS != nil {
		r := response.New()
		if err = self.PreCSS.Render(r); err != nil {
			return err
		}
		templateContext.PreCSS = r.String()
	}
	if self.CSS != nil {
		templateContext.CSS = self.CSS.URL(response.Request.Params...)
	} else {
		templateContext.CSS = Config.Page.DefaultCSS
	}
	if self.PostCSS != nil {
		r := response.New()
		if err = self.PostCSS.Render(r); err != nil {
			return err
		}
		templateContext.PostCSS = r.String()
	}

	headScripts := self.HeadScripts
	if headScripts == nil {
		headScripts = Config.Page.DefaultHeadScripts
	}
	if headScripts != nil {
		r := response.New()
		if err = headScripts.Render(r); err != nil {
			return err
		}
		templateContext.HeadScripts = r.String()
	}

	scripts := self.Scripts
	if scripts == nil {
		scripts = Config.Page.DefaultScripts
	}
	if scripts != nil {
		r := response.New()
		if err = scripts.Render(r); err != nil {
			return err
		}
		if Config.Page.PostScripts != nil {
			if err = Config.Page.PostScripts.Render(r); err != nil {
				return err
			}
		}
		templateContext.Scripts = r.String()
	}

	if self.Content != nil {
		r := response.New()
		err = self.Content.Render(r)
		if err != nil {
			return err
		}
		templateContext.Content = r.String()
	}

	// Get dynamic style and scripts after self.Content.Render()
	// because they are added in Render()
	templateContext.DynamicStyle = context.dynamicStyle.String()
	templateContext.DynamicHeadScripts = context.dynamicHeadScripts.String()
	templateContext.DynamicScripts = context.dynamicScripts.String()

	self.Template.GetContext = TemplateContext(templateContext)
	return self.Template.Render(response)
}
