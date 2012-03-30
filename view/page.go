package view

import (
	"bytes"
	"github.com/ungerik/go-start/templatesystem"
	"github.com/ungerik/go-start/utils"
	"html"
	"io"
)

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
		writer.Write([]byte(url.URL(context)))
		writer.Write([]byte("'>\n"))
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////////
// Page

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

Most HTML head and script specific text is written by PageWriteFunc functions
that receive the request context as an argument.
That way the content can be created dynamically.

Wrapper functions for static content are provided for convenience.
See functions under PageWriteFunc below.

Example:

	&Page{WriteTitle: func(context *Context, writer io.Writer) (err error) {
		writer.Write([]byte("Could be a dynamic title"))
		return nil
	}}

	&Page{WriteTitle: PageTitle("Static Title")}


To avoid reading the same data multiple times from the database in PageWriteFunc
or dynamic views in the content structure, OnPreRender can be used to
query and set page wide data only once at the request context.Data.

Example:

	&Page{
		OnPreRender: func(page *Page, context *Context) (err error) {
			context.Data = &MyPerPageData{SomeText: "Hello World!"}
		},
		Content: NewDynamicView(
			func(context *Context) (view View, err error) {
				myPerPageData := context.Data.(*MyPerPageData)
				return HTML(myPerPageData.SomeText), nil
			},
		),
	}

*/
type Page struct {
	Template

	// Called before any other function when rendering the page
	OnPreRender func(page *Page, context *Context) (err error)

	// Writes the head title tag (HTML escaped)
	WriteTitle PageWriteFunc

	// Writes the head meta description tag (HTML escaped)
	WriteMetaDescription PageWriteFunc

	// Content of the head meta viewport tag,
	// Config.Page.DefaultMetaViewport will be used if ""
	MetaViewport string

	// Write additional HTML head content
	WriteHeader PageWriteFunc

	// Write head content before the stylesheet link
	WritePreCSS PageWriteFunc

	// stylesheet link URL
	CSS URL

	// Write head content after the stylesheet link
	WritePostCSS PageWriteFunc

	// HTML body content. Will be wrapped by a div with class="container"
	Content View

	// Write scripts after body content
	WriteScripts PageWriteFunc

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
func (self *Page) URL(context *Context, args ...string) string {
	path := StringURL(self.path).URL(context, args...)
	return "http://" + context.Request.Host + path
}

// Implements the LinkModel interface
func (self *Page) LinkContent(context *Context) View {
	return HTML(self.LinkTitle(context))
}

// Implements the LinkModel interface
func (self *Page) LinkTitle(context *Context) string {
	if self.WriteTitle == nil {
		return ""
	}
	var buf bytes.Buffer
	err := self.WriteTitle(context, &buf)
	if err != nil {
		//return err.String()
		panic(err)
	}
	return buf.String()
}

// Implements the LinkModel interface
func (self *Page) LinkRel(context *Context) string {
	return ""
}

func (self *Page) Render(context *Context, writer *utils.XMLWriter) (err error) {
	if self.OnPreRender != nil {
		err = self.OnPreRender(self, context)
		if err != nil {
			return err
		}
	}

	var templateContext struct {
		Title             string
		MetaDescription   string
		MetaViewport      string
		Header            string
		PreCSS            string
		CSS               string
		PostCSS           string
		Scripts           string
		Favicon16x16URL   string
		Favicon57x57URL   string
		Favicon72x72URL   string
		Favicon114x114URL string
		Favicon129x129URL string
		Content           string
	}

	if self.WriteTitle != nil {
		var buf bytes.Buffer
		err := self.WriteTitle(context, &buf)
		if err != nil {
			return err
		}
		templateContext.Title = html.EscapeString(buf.String())
	}
	if self.WriteMetaDescription != nil {
		var buf bytes.Buffer
		err := self.WriteMetaDescription(context, &buf)
		if err != nil {
			return err
		}
		templateContext.MetaDescription = html.EscapeString(buf.String())
	}

	metaViewport := self.MetaViewport
	if metaViewport == "" {
		metaViewport = Config.Page.DefaultMetaViewport
	}
	templateContext.MetaViewport = metaViewport

	writeHeader := self.WriteHeader
	if writeHeader == nil {
		writeHeader = Config.Page.DefaultWriteHeader
	}
	if writeHeader != nil {
		var buf bytes.Buffer
		err := writeHeader(context, &buf)
		if err != nil {
			return err
		}
		templateContext.Header = buf.String()
	}

	//templateContext.Meta = self.Meta
	templateContext.Favicon16x16URL = self.Favicon16x16URL
	templateContext.Favicon57x57URL = self.Favicon57x57URL
	templateContext.Favicon72x72URL = self.Favicon72x72URL
	templateContext.Favicon114x114URL = self.Favicon114x114URL
	templateContext.Favicon129x129URL = self.Favicon129x129URL

	if self.WritePreCSS != nil {
		var buf bytes.Buffer
		if err = self.WritePreCSS(context, &buf); err != nil {
			return err
		}
		templateContext.PreCSS = buf.String()
	}
	if self.CSS != nil {
		templateContext.CSS = self.CSS.URL(context)
	} else {
		templateContext.CSS = Config.Page.DefaultCSS
	}
	if self.WritePostCSS != nil {
		var buf bytes.Buffer
		if err = self.WritePostCSS(context, &buf); err != nil {
			return err
		}
		templateContext.PostCSS = buf.String()
	}

	writeScripts := self.WriteScripts
	if writeScripts == nil {
		writeScripts = Config.Page.DefaultWriteScripts
	}
	if writeScripts != nil {
		var buf bytes.Buffer
		if err = writeScripts(context, &buf); err != nil {
			return err
		}
		if Config.Page.PostWriteScripts != nil {
			if err = Config.Page.PostWriteScripts(context, &buf); err != nil {
				return err
			}
		}
		templateContext.Scripts = buf.String()
	}

	contentHtml := utils.NewXMLBuffer()
	if self.Content != nil {
		err = self.Content.Render(context, &contentHtml.XMLWriter)
		if err != nil {
			return err
		}
	}
	templateContext.Content = contentHtml.String()

	self.Template.GetContext = TemplateContext(templateContext)
	return self.Template.Render(context, writer)
}
