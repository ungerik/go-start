package view

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/templatesystem"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
	"path/filepath"
)

type PageConfiguration struct {
	Template            string
	DefaultWriteHeader  PageWriteFunc // will be called after WriteTitle
	DefaultCSS          string
	DefaultMetaViewport string
	DefaultWriteScripts PageWriteFunc // will be called if Page.WriteScripts is nil
	PostWriteScripts    PageWriteFunc // will always be called after Page.WriteScripts
	DefaultAuth         Authenticator // Will be used for pages with Page.NeedsAuth == true
}

type Configuration struct {
	TemplateSystem            templatesystem.Implementation
	Page                      PageConfiguration
	BaseDirs                  []string
	StaticDirs                []string
	TemplateDirs              []string
	RedirectSubdomains        []string // Exapmle: "www"
	BaseURL                   string
	SiteName                  string
	CookieSecret              string
	SessionTracker            SessionTracker
	SessionDataStore          SessionDataStore
	GlobalAuth                Authenticator // Will allways be used before all other authenticators
	FallbackAuth              Authenticator // Will be used when no other authenticator is defined for the view
	LoginSignupPage           **Page
	Middlewares               []Middleware
	NumFieldRepeatFormMessage int
	Debug                     bool
	FormErrorMessageClass     string
	FormSuccessMessageClass   string
	DebugPrintPaths           bool
}

// Config holds the configuration of the view package.
var Config Configuration = Configuration{
	TemplateSystem: &templatesystem.Mustache{},
	Page: PageConfiguration{
		Template:            "html5boilerplate.html",
		DefaultCSS:          "/style.css",
		DefaultMetaViewport: "width=device-width",
		//DefaultWriteScripts: JQuery,
	},
	BaseDirs:                  []string{"."},
	StaticDirs:                []string{"static"},    // every StaticDir will be appended to every BaseDir to search for static files
	TemplateDirs:              []string{"templates"}, // every TemplateDir will be appended to every BaseDir to search for template files
	SessionTracker:            &CookieSessionTracker{},
	SessionDataStore:          NewCookieSessionDataStore(),
	NumFieldRepeatFormMessage: 6,
	FormErrorMessageClass:     "error",
	FormSuccessMessageClass:   "success",
}

// Init updates Config with the site-name, cookie secret and base directories used for static and template file search.
// For every directory of baseDirs, Config.StaticDirs are appended to create search paths for static files
// and Config.TemplateDirs are appended to search for template files.
func Init(siteName, cookieSecret string, baseDirs ...string) (err error) {
	Config.SiteName = siteName
	Config.CookieSecret = cookieSecret
	if len(baseDirs) > 0 {
		for i, dir := range baseDirs {
			dir, err = filepath.Abs(dir)
			if err != nil {
				return err
			}
			if !utils.DirExists(dir) {
				return errs.Format("BaseDir does not exist: %s", dir)
			}
			baseDirs[i] = dir
		}

		Config.BaseDirs = baseDirs
	}
	return nil
}

func Close() {
	web.Close()
}
