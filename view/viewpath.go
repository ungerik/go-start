package view

import (
	"fmt"
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	runtime_debug "runtime/debug"
	"strings"
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
				self.NoAuth = IndirectURL(Config.LoginSignupPage)
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

	if Config.Debug.PrintPaths {
		debug.Print(path)
	}

	//debug.Print(path)
	self.View.Init(self.View)
	if viewWithURL, ok := self.View.(ViewWithURL); ok {
		viewWithURL.SetPath(path)
	}

	htmlFunc := func(webContext *web.Context, args ...string) string {

		// See: http://groups.google.com/group/golang-nuts/browse_thread/thread/ab1971bb9459025d
		// Slows down the memory leak on 32 bit systems a little bit:
		defer func() {
			go runtime.GC()
		}()

		response := newResponse(webContext, self.View, args)

		for _, subdomain := range Config.RedirectSubdomains {
			if len(subdomain) > 0 {
				if subdomain[len(subdomain)-1] != '.' {
					subdomain += "."
				}
				host := response.Request.Host
				if strings.Index(host, subdomain) == 0 {
					host = host[len(subdomain):]
					url := "http://" + host + response.Request.URL.Path
					response.RedirectPermanently301(url)
					return ""
				}
			}
		}

		handleErr := func(err error) string {
			switch err.(type) {
			case NotFound:
				response.NotFound404(err.Error())
			case Redirect:
				if Config.Debug.PrintRedirects {
					fmt.Printf("%d Redirect: %s\n", http.StatusFound, err.Error())
				}
				response.RedirectTemporary302(err.Error())
			case PermanentRedirect:
				if Config.Debug.PrintRedirects {
					fmt.Printf("%d Permanent Redirect: %s\n", http.StatusMovedPermanently, err.Error())
				}
				response.RedirectPermanently301(err.Error())
			case Forbidden:
				response.Forbidden403(err.Error())
			default:
				fmt.Println(err.Error())
				runtime_debug.PrintStack()
				msg := err.Error()
				if Config.Debug.Mode {
					msg += string(runtime_debug.Stack())
				}
				response.Abort(http.StatusInternalServerError, msg)
			}
			return ""
		}

		handleNoAuth := func(err error) string {
			switch {
			case err != nil:
				return handleErr(err)
			case self.NoAuth != nil:
				from := url.QueryEscape(response.Request.RequestURI)
				to := self.NoAuth.URL(response) + "?from=" + from
				return handleErr(Redirect(to))
			}
			return handleErr(Forbidden("403 Forbidden: authentication required"))
		}

		if Config.OnPreAuth != nil {
			if err := Config.OnPreAuth(response); err != nil {
				return handleErr(err)
			}
		}

		if Config.GlobalAuth != nil {
			if ok, err := Config.GlobalAuth.Authenticate(response); !ok {
				return handleNoAuth(err)
			}
		}

		if self.Auth == nil {
			self.Auth = Config.FallbackAuth
		}
		if self.Auth != nil {
			if ok, err := self.Auth.Authenticate(response); !ok {
				return handleNoAuth(err)
			}
		}

		//		numMiddlewares := len(Config.Middlewares)
		//		for i := 0; i < numMiddlewares; i++ {
		//			if abort := Config.Middlewares[i].PreRender(context); abort {
		//				return ""
		//			}
		//		}

		err := self.View.Render(response)

		//		for i := numMiddlewares - 1; i >= 0; i-- {
		//			html, err = Config.Middlewares[i].PostRender(context, html, err)
		//		}

		if err != nil {
			return handleErr(err)
		}
		return response.String()
	}

	web.Get(path, htmlFunc)
	web.Post(path, htmlFunc)

	for i := range self.Sub {
		self.Sub[i].initAndRegisterViewsRecursive(path)
	}
}
