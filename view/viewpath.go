package view

import (
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/utils"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"github.com/ungerik/web.go"
)

const PathFragmentPattern = "([a-zA-Z0-9_\\-\\.]+)"

var PathFragmentRegexp *regexp.Regexp = regexp.MustCompile(PathFragmentPattern)

///////////////////////////////////////////////////////////////////////////////
// ViewPath

// ViewPath holds all data necessary to define the URL path of a view,
// including the number of arguments parsed from the URL path,
// an Authenticator and sub paths.
type ViewPath struct {
	Name   string
	Args   int
	View   View
	Auth   Authenticator
	NoAuth URL
	Sub    []ViewPath // Only allowed when View is a Page or nil
}

// Pages and nil views will be registered with a trailing slash at their path
// and a permanent redirect from the path without trailing slash
// Nil views will be registered as NotFound404
// parentPath always ends with a slash
func (self *ViewPath) initAndRegisterViewsRecursive(parentPath string) {
	if parentPath == "" || parentPath[len(parentPath)-1] != '/' {
		panic("Parent path must end with a slash: " + parentPath)
	}
	if parentPath != "/" && self.Name == "" && self.Args == 0 {
		panic("Sub path of " + parentPath + " with no Name and Args")
	}
	if self.Name != "" && !PathFragmentRegexp.MatchString(self.Name) {
		panic("Invalid characters in view.ViewPath.Name: " + self.Name)
	}
	if self.View != nil && utils.IsDeepNil(self.View) {
		panic("Nil value wrapped with non nil view.View under parentPath: " + parentPath)
	}

	addSlash := self.Args > 0
	if self.View == nil {
		addSlash = true
		self.View = &NotFoundView{Message: "Invalid URL"}
	} else if _, isPage := self.View.(*Page); isPage {
		addSlash = true
		if self.Auth != nil {
			if self.NoAuth == nil && Config.LoginSignupPage != nil && *Config.LoginSignupPage != nil {
				self.NoAuth = &IndirectPageURL{Config.LoginSignupPage}
			}
		}
	}

	path := parentPath + self.Name
	if self.Args > 0 {
		if self.Name != "" {
			path += "/"
		}
		for i := 0; i < self.Args; i++ {
			path += PathFragmentPattern + "/"
		}
	}

	if addSlash {
		if path[len(path)-1] != '/' {
			path += "/"
		}
		if path != "/" {
			web.Get(path[:len(path)-1], func(webContext *web.Context, args ...string) string {
				webContext.Redirect(http.StatusMovedPermanently, webContext.Request.RequestURI+"/")
				return ""
			})
		}
	}

	if self.Args < 0 {
		panic("Negative Args at " + path)
	}
	if _, pathExists := viewsByPath[path]; pathExists {
		panic("View with path '" + path + "' already registered")
	}
	viewsByPath[path] = self.View

	if Config.DebugPrintPaths {
		debug.Print(path)
	}

	//debug.Print(path)
	self.View.Init(self.View)
	if viewWithURL, ok := self.View.(ViewWithURL); ok {
		viewWithURL.SetPath(path)
	}

	htmlFunc := func(webContext *web.Context, args ...string) string {
		context := NewContext(webContext, self.View, args)

		for _, subdomain := range Config.RedirectSubdomains {
			if len(subdomain) > 0 {
				if subdomain[len(subdomain)-1] != '.' {
					subdomain += "."
				}
				host := context.Request.Host
				if strings.Index(host, subdomain) == 0 {
					host = host[len(subdomain):]
					url := "http://" + host + context.Request.URL.Path
					context.Redirect(http.StatusMovedPermanently, url)
					return ""
				}
			}
		}

		handleErr := func(err error) string {
			switch err.(type) {
			case NotFound:
				context.NotFound(err.Error())
			case Redirect:
				context.Redirect(http.StatusFound, err.Error())
			case PermanentRedirect:
				context.Redirect(http.StatusMovedPermanently, err.Error())
			case Forbidden:
				context.Abort(http.StatusForbidden, err.Error())
			default:
				context.Abort(http.StatusInternalServerError, err.Error())
			}
			return ""
		}

		handleNoAuth := func(err error) string {
			switch {
			case err != nil:
				return handleErr(err)
			case self.NoAuth != nil:
				from := url.QueryEscape(context.Request.RequestURI)
				to := self.NoAuth.URL(context) + "?from=" + from
				return handleErr(Redirect(to))
			}
			return handleErr(Forbidden("403 Forbidden: authentication required"))
		}

		if Config.GlobalAuth != nil {
			if ok, err := Config.GlobalAuth.Authenticate(context); !ok {
				return handleNoAuth(err)
			}
		}

		if self.Auth == nil {
			self.Auth = Config.FallbackAuth
		}
		if self.Auth != nil {
			if ok, err := self.Auth.Authenticate(context); !ok {
				return handleNoAuth(err)
			}
		}

		//		numMiddlewares := len(Config.Middlewares)
		//		for i := 0; i < numMiddlewares; i++ {
		//			if abort := Config.Middlewares[i].PreRender(context); abort {
		//				return ""
		//			}
		//		}

		xmlBuffer := utils.NewXMLBuffer()

		err := self.View.Render(context, &xmlBuffer.XMLWriter)

		//		for i := numMiddlewares - 1; i >= 0; i-- {
		//			html, err = Config.Middlewares[i].PostRender(context, html, err)
		//		}

		if err != nil {
			return handleErr(err)
		}
		return xmlBuffer.String()
	}

	web.Get(path, htmlFunc)
	web.Post(path, htmlFunc)

	for i := range self.Sub {
		self.Sub[i].initAndRegisterViewsRecursive(path)
	}
}
