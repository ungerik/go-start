package view

import (
	"html"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/templatesystem"
)

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
		OnPreRender: func(page *Page, ctx *Context) (err error) {
			context.Data = &MyPerPageData{SomeText: "Hello World!"}
		},
		Content: DynamicView(
			func(ctx *Context) (view View, err error) {
				myPerPageData := context.Data.(*MyPerPageData)
				return HTML(myPerPageData.SomeText), nil
			},
		),
	}

*/
type Page struct {
	Template

	// Called before any other function when rendering the page
	OnPreRender func(page *Page, ctx *Context) (err error)

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
func (self *Page) URL(ctx *Context) string {
	if self == nil {
		panic("view.Page is nil, potential circular initialization dependency, consider using addresses of page variables")
	}
	return StringURL(self.path).URL(ctx)
}

// Implements the LinkModel interface
func (self *Page) LinkContent(ctx *Context) View {
	if self == nil {
		panic("view.Page is nil, potential circular initialization dependency, consider using addresses of page variables")
	}
	return HTML(self.LinkTitle(ctx))
}

// Implements the LinkModel interface
func (self *Page) LinkTitle(ctx *Context) string {
	if self == nil {
		panic("view.Page is nil, potential circular initialization dependency, consider using addresses of page variables")
	}
	if self.Title == nil {
		return ""
	}
	ctx.Response.PushBody()
	err := self.Title.Render(ctx)
	if err != nil {
		panic(err)
	}
	return ctx.Response.PopBodyString()
}

// Implements the LinkModel interface
func (self *Page) LinkRel(ctx *Context) string {
	return ""
}

func (self *Page) Render(ctx *Context) (err error) {
	if self.OnPreRender != nil {
		err = self.OnPreRender(self, ctx)
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
		ctx.Response.PushBody()
		err := self.Title.Render(ctx)
		if err != nil {
			return err
		}
		templateContext.Title = html.EscapeString(ctx.Response.PopBodyString())
	}
	if self.MetaDescription != nil {
		ctx.Response.PushBody()
		err := self.MetaDescription.Render(ctx)
		if err != nil {
			return err
		}
		templateContext.MetaDescription = html.EscapeString(ctx.Response.PopBodyString())
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
		ctx.Response.PushBody()
		err := additionalHead.Render(ctx)
		if err != nil {
			return err
		}
		templateContext.Head = ctx.Response.PopBodyString()
	}

	//templateContext.Meta = self.Meta
	templateContext.Favicon16x16URL = self.Favicon16x16URL
	templateContext.Favicon57x57URL = self.Favicon57x57URL
	templateContext.Favicon72x72URL = self.Favicon72x72URL
	templateContext.Favicon114x114URL = self.Favicon114x114URL
	templateContext.Favicon129x129URL = self.Favicon129x129URL

	if self.PreCSS != nil {
		ctx.Response.PushBody()
		if err = self.PreCSS.Render(ctx); err != nil {
			return err
		}
		templateContext.PreCSS = ctx.Response.PopBodyString()
	}
	if self.CSS != nil {
		templateContext.CSS = self.CSS.URL(ctx)
	} else {
		templateContext.CSS = Config.Page.DefaultCSS
	}
	if self.PostCSS != nil {
		ctx.Response.PushBody()
		if err = self.PostCSS.Render(ctx); err != nil {
			return err
		}
		templateContext.PostCSS = ctx.Response.PopBodyString()
	}

	headScripts := self.HeadScripts
	if headScripts == nil {
		headScripts = Config.Page.DefaultHeadScripts
	}
	if headScripts != nil {
		ctx.Response.PushBody()
		if err = headScripts.Render(ctx); err != nil {
			return err
		}
		templateContext.HeadScripts = ctx.Response.PopBodyString()
	}

	scripts := self.Scripts
	if scripts == nil {
		scripts = Config.Page.DefaultScripts
	}
	if scripts != nil {
		ctx.Response.PushBody()
		if err = scripts.Render(ctx); err != nil {
			return err
		}
		if Config.Page.PostScripts != nil {
			if err = Config.Page.PostScripts.Render(ctx); err != nil {
				return err
			}
		}
		templateContext.Scripts = ctx.Response.PopBodyString()
	}

	if self.Content != nil {
		ctx.Response.PushBody()
		err = self.Content.Render(ctx)
		if err != nil {
			return err
		}
		templateContext.Content = ctx.Response.PopBodyString()
	}

	// Get dynamic style and scripts after self.Content.Render()
	// because they are added in Render()
	templateContext.DynamicStyle = ctx.Response.dynamicStyle.String()
	templateContext.DynamicHeadScripts = ctx.Response.dynamicHeadScripts.String()
	templateContext.DynamicScripts = ctx.Response.dynamicScripts.String()

	self.Template.GetContext = TemplateContext(templateContext)
	return self.Template.Render(ctx)
}
