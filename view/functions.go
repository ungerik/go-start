package view

import (
	"io"
	"log"
	"strconv"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/templatesystem"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/web.go"

//	"github.com/ungerik/go-start/debug"
)

var viewIdChan = make(chan string, 16)
var viewIdCounter int64

func init() {
	go func() {
		for {
			viewIdCounter++
			// Use base32 encoding for ids to make them shorter
			viewIdChan <- strconv.FormatInt(viewIdCounter, 32)
		}
	}()
}

// var viewsByID map[string]View = map[string]View{}
var viewsByPath = map[string]View{}

func NewViewID(view View) (id string) {
	id = <-viewIdChan
	// viewsByID[id] = view
	return id
}

// func DeleteViewID(id string) {
// 	if _, exists := viewsByID[id]; !exists {
// 		panic("View ID '" + id + "' does not exist")
// 	}
// 	delete(viewsByID, id)
// }

func FindStaticFile(filename string) (filePath string, found bool, modifiedTime int64) {
	return utils.FindFile2ModifiedTime(Config.BaseDirs, Config.StaticDirs, filename)
}

func FindTemplateFile(filename string) (filePath string, found bool, modifiedTime int64) {
	return utils.FindFile2ModifiedTime(Config.BaseDirs, Config.TemplateDirs, filename)
}

//func ViewChanged(view View) {
//}

// RunServer starts a webserver with the given paths.
// If paths is nil, only static files will be served.
func RunServer(paths *ViewPath) {
	if !Config.initialized {
		err := Config.Init()
		if err != nil {
			panic(err)
		}
	}
	addr := Config.ListenAndServeAt
	if !Config.IsProductionServer && Config.Debug.ListenAndServeAt != "" {
		addr = Config.Debug.ListenAndServeAt
	}
	RunServerAddr(addr, paths)
}

// RunServerAddr starts a webserver with the given paths and address.
// If paths is nil, only static files will be served.
func RunServerAddr(addr string, paths *ViewPath) {
	if !Config.initialized {
		err := Config.Init()
		if err != nil {
			panic(err)
		}
	}
	log.Print("view.Config.IsProductionServer = ", Config.IsProductionServer)
	log.Print("view.Config.Debug.Mode = ", Config.Debug.Mode)

	if paths != nil {
		paths.initAndRegisterViewsRecursive("/")
	}

	web.Config.StaticDirs = utils.CombineDirs(Config.BaseDirs, Config.StaticDirs)
	web.Config.RecoverPanic = Config.Debug.Mode
	web.Config.CookieSecret = Config.CookieSecret

	web.Run(addr)
}

func RenderTemplate(filename string, out io.Writer, context interface{}) (err error) {
	filePath, found, _ := FindTemplateFile(filename)
	if !found {
		return errs.Format("Template file not found: %s", filename)
	}

	var templ templatesystem.Template
	templ, err = Config.TemplateSystem.ParseFile(filePath)
	if err != nil {
		return
	}

	// context = append(context, Config)
	return templ.Render(out, context)
}

func RenderTemplateString(tmplString string, name string, out io.Writer, context interface{}) (err error) {
	var templ templatesystem.Template
	templ, err = Config.TemplateSystem.ParseString(tmplString, name)
	if err != nil {
		return
	}

	// context = append(context, Config)
	return templ.Render(out, context)
}

func RenderChildViewsHTML(parent View, ctx *Context) (err error) {
	parent.IterateChildren(func(parent View, child View) (next bool) {
		if child != nil {
			err = child.Render(ctx)
			if err != nil {
				return false
			}
		}
		return true
	})
	return err
}

