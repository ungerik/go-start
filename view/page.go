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
func (self *Page) URL(response *Response) string {
	return StringURL(self.path).URL(response)
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
	oldBody := response.SwitchBody(nil)
	err := self.Title.Render(response)
	if err != nil {
		panic(err)
	}
	return string(response.SwitchBody(oldBody))
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

	originalBody := response.SwitchBody(nil)

	if self.Title != nil {
		err := self.Title.Render(response)
		if err != nil {
			return err
		}
		templateContext.Title = html.EscapeString(string(response.SwitchBody(nil)))
	}
	if self.MetaDescription != nil {
		err := self.MetaDescription.Render(response)
		if err != nil {
			return err
		}
		templateContext.MetaDescription = html.EscapeString(string(response.SwitchBody(nil)))
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
		err := additionalHead.Render(response)
		if err != nil {
			return err
		}
		templateContext.Head = string(response.SwitchBody(nil))
	}

	//templateContext.Meta = self.Meta
	templateContext.Favicon16x16URL = self.Favicon16x16URL
	templateContext.Favicon57x57URL = self.Favicon57x57URL
	templateContext.Favicon72x72URL = self.Favicon72x72URL
	templateContext.Favicon114x114URL = self.Favicon114x114URL
	templateContext.Favicon129x129URL = self.Favicon129x129URL

	if self.PreCSS != nil {
		if err = self.PreCSS.Render(response); err != nil {
			return err
		}
		templateContext.PreCSS = string(response.SwitchBody(nil))
	}
	if self.CSS != nil {
		templateContext.CSS = self.CSS.URL(response)
	} else {
		templateContext.CSS = Config.Page.DefaultCSS
	}
	if self.PostCSS != nil {
		if err = self.PostCSS.Render(response); err != nil {
			return err
		}
		templateContext.PostCSS = string(response.SwitchBody(nil))
	}

	headScripts := self.HeadScripts
	if headScripts == nil {
		headScripts = Config.Page.DefaultHeadScripts
	}
	if headScripts != nil {
		if err = headScripts.Render(response); err != nil {
			return err
		}
		templateContext.HeadScripts = string(response.SwitchBody(nil))
	}

	scripts := self.Scripts
	if scripts == nil {
		scripts = Config.Page.DefaultScripts
	}
	if scripts != nil {
		if err = scripts.Render(response); err != nil {
			return err
		}
		if Config.Page.PostScripts != nil {
			if err = Config.Page.PostScripts.Render(response); err != nil {
				return err
			}
		}
		templateContext.Scripts = string(response.SwitchBody(nil))
	}

	if self.Content != nil {
		err = self.Content.Render(response)
		if err != nil {
			return err
		}
		templateContext.Content = string(response.SwitchBody(nil))
	}

	// Get dynamic style and scripts after self.Content.Render()
	// because they are added in Render()
	templateContext.DynamicStyle = response.dynamicStyle.String()
	templateContext.DynamicHeadScripts = response.dynamicHeadScripts.String()
	templateContext.DynamicScripts = response.dynamicScripts.String()

	response.SwitchBody(originalBody)

	self.Template.GetContext = TemplateContext(templateContext)
	return self.Template.Render(response)
}
