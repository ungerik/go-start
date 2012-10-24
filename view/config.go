package view

import (
	"net"
	"path/filepath"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/templatesystem"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"
)

var Config = Configuration{
	ListenAndServeAt: "0.0.0.0:80",
	TemplateSystem:   &templatesystem.Mustache{},
	Page: PageConfiguration{
		Template:            "html5boilerplate.html",
		DefaultCSS:          "/style.css",
		DefaultMetaViewport: "width=device-width",
		//DefaultWriteScripts: JQuery,
	},
	Form: FormConfiguration{
		DefaultLayout: &StandardFormLayout{
			DefaultInputSize:      80,
			DefaultTableInputSize: 20,
		},
		DefaultCSRFProtector:            nil,
		DefaultSubmitButtonText:         "Save",
		DefaultErrorMessageClass:        "error",
		DefaultSuccessMessageClass:      "success",
		DefaultFieldDescriptionClass:    "description",
		DefaultRequiredMarker:           HTML("<span class='required'>*</span>"),
		GeneralErrorMessageOnFieldError: "This form has errors",
		DefaultFieldControllers: FormFieldControllers{
			ModelStringController{},
			ModelTextController{},
			ModelRichTextController{},
			ModelUrlController{},
			ModelEmailController{},
			ModelPasswordController{},
			ModelIntController{},
			ModelFloatController{},
			ModelPhoneController{},
			ModelBoolController{},
			ModelChoiceController{},
			ModelMultipleChoiceController{},
			ModelDynamicChoiceController{},
			ModelDateController{},
			ModelDateTimeController{},
			ModelFileController{},
			ModelBlobController{},
		},
	},
	RichText: RichTextConfiguration{
		DefaultToolbar: SetDefaultToolbar(),
		ToolbarCSS:     "/css/wysihtml5/toolbar-stylesheet.css",
		EditorCSS:      "",
		UseGlobalCSS:   false,
	},
	LabeledModelViewLabelClass: "labeled-model-view-label",
	LabeledModelViewValueClass: "labeled-model-view-value",
	BaseDirs:                   []string{"."},
	StaticDirs:                 []string{"static"},    // every StaticDir will be appended to every BaseDir to search for static files
	TemplateDirs:               []string{"templates"}, // every TemplateDir will be appended to every BaseDir to search for template files
	SessionTracker:             &CookieSessionTracker{},
	SessionDataStore:           NewCookieSessionDataStore(),
	NamedAuthenticators:        make(map[string]Authenticator),
}

var StructTagKey = "view"

type Configuration struct {
	initialized                bool
	ListenAndServeAt           string
	IsProductionServer         bool // IsProductionServer will be set to true if localhost resolves to one of ProductionServerIPs
	ProductionServerIPs        []string
	TemplateSystem             templatesystem.Implementation
	Page                       PageConfiguration
	Form                       FormConfiguration
	RichText                   RichTextConfiguration
	LabeledModelViewLabelClass string
	LabeledModelViewValueClass string
	DisableCachedViews         bool
	BaseDirs                   []string
	StaticDirs                 []string
	TemplateDirs               []string
	RedirectSubdomains         []string // Exapmle: "www"
	SiteName                   string
	CookieSecret               string
	SessionTracker             SessionTracker
	SessionDataStore           SessionDataStore
	OnPreAuth                  func(ctx *Context) error
	GlobalAuth                 Authenticator // Will allways be used before all other authenticators
	FallbackAuth               Authenticator // Will be used when no other authenticator is defined for the view
	NamedAuthenticators        map[string]Authenticator
	LoginSignupPage            **Page
	// Middlewares               []Middleware
	Debug struct {
		ListenAndServeAt string
		Mode             bool // Will be set to true if IsProductionServer is false
		LogPaths         bool
		LogRedirects     bool
	}
}

func (self *Configuration) Name() string {
	return "view"
}

func (self *Configuration) Init() error {
	if self.initialized {
		panic("view.Config already initialized")
	}

	if !self.IsProductionServer {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				ip := ipNet.IP.String()
				for _, prodIP := range Config.ProductionServerIPs {
					if ip == prodIP {
						self.IsProductionServer = true
						break
					}
				}
			}
		}
	}

	if !self.IsProductionServer {
		self.Debug.Mode = true
	}

	// Check if dir exists and make it absolute
	for i := range Config.BaseDirs {
		dir, err := filepath.Abs(Config.BaseDirs[i])
		if err != nil {
			return err
		}
		if !utils.DirExists(dir) {
			return errs.Format("BaseDir does not exist: %s", dir)
		}
		Config.BaseDirs[i] = dir
	}

	self.initialized = true
	return nil
}

func (self *Configuration) Close() error {
	web.Close()
	return nil
}

type PageConfiguration struct {
	Template              string
	DefaultAdditionalHead Renderer // will be called after WriteTitle
	DefaultCSS            string
	DefaultMetaViewport   string
	DefaultHeadScripts    Renderer      // write scripts as last element of the HTML head
	DefaultScripts        Renderer      // will be called if Page.WriteScripts is nil
	PostScripts           Renderer      // will always be called after Page.WriteScripts
	DefaultAuth           Authenticator // Will be used for pages with Page.NeedsAuth == true
}

type FormConfiguration struct {
	DefaultLayout                   FormLayout
	DefaultCSRFProtector            CSRFProtector
	DefaultErrorMessageClass        string
	DefaultSuccessMessageClass      string
	DefaultSubmitButtonClass        string
	DefaultFieldDescriptionClass    string
	StandardFormLayoutDivClass      string
	DefaultSubmitButtonText         string
	GeneralErrorMessageOnFieldError string
	DefaultRequiredMarker           View
	DefaultFieldControllers         FormFieldControllers
}

type RichTextConfiguration struct {
	DefaultToolbar string
	ToolbarCSS     string
	EditorCSS      string
	UseGlobalCSS   bool
}

func (self *RichTextConfiguration) SetStylesheet(url string) {
	if url != "" {
		Config.RichText.ToolbarCSS = url
	} else {
		Config.RichText.UseGlobalCSS = true
	}
}

// // Init updates Config with the site-name, cookie secret and base directories used
// // for static and template file search.
// // For every directory of baseDirs, Config.StaticDirs are appended to create
// // search paths for static files and Config.TemplateDirs are appended
// // to search for template files.
// func Init(siteName, cookieSecret string, baseDirs ...string) (err error) {
// 	Config.SiteName = siteName
// 	Config.CookieSecret = cookieSecret
// 	if len(baseDirs) > 0 {
// 		Config.BaseDirs = baseDirs
// 	}
// 	return Config.Init()
// }
