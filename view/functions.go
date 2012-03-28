package view

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/templatesystem"
	"github.com/ungerik/go-start/utils"
	"github.com/ungerik/goconfig"
	"io"
	"strconv"
	"github.com/ungerik/web.go"

//	"github.com/ungerik/go-start/debug"
)

var viewsByID map[string]View = map[string]View{}
var viewsByPath map[string]View = map[string]View{}

var viewIdCounter uint64

func NewViewID(view View) (id string) {
	viewIdCounter++
	id = strconv.FormatUint(viewIdCounter, 32)
	viewsByID[id] = view
	return id
}

func DeleteViewID(id string) {
	if _, exists := viewsByID[id]; !exists {
		panic("View ID '" + id + "' does not exist")
	}
	delete(viewsByID, id)
}

func FindStaticFile(filename string) (filePath string, found bool, modifiedTime int64) {
	return utils.FindFile2ModifiedTime(Config.BaseDirs, Config.StaticDirs, filename)
}

func FindTemplateFile(filename string) (filePath string, found bool, modifiedTime int64) {
	return utils.FindFile2ModifiedTime(Config.BaseDirs, Config.TemplateDirs, filename)
}

//func ViewChanged(view View) {
//}

func Run(paths *ViewPath, address, baseURL string, recoverPanic bool) {
	Config.BaseURL = baseURL

	paths.initAndRegisterViewsRecursive("/")

	staticDirs := utils.CombineDirs(Config.BaseDirs, Config.StaticDirs)
	web.Config.StaticDirs = staticDirs
	web.Config.RecoverPanic = recoverPanic
	web.Config.CookieSecret = Config.CookieSecret
	web.Run(address)
}

func RunConfigFile(paths *ViewPath, filename string) {
	config, err := config.ReadDefault(filename)
	if err != nil {
		panic(err)
	}

	var host string
	host, err = config.String("server", "host")
	if err != nil {
		panic(err)
	}

	var port string
	port, err = config.String("server", "port")
	if err != nil {
		panic(err)
	}

	var baseURL string
	baseURL, err = config.String("server", "url")
	if err != nil {
		panic(err)
	}

	var recoverPanic bool
	recoverPanic, _ = config.Bool("server", "recoverPanic")

	Run(paths, host+":"+port, baseURL, recoverPanic)
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

func RenderChildViewsHTML(parent View, context *Context, writer *utils.XMLWriter) (err error) {
	parent.IterateChildren(func(parent View, child View) (next bool) {
		if child != nil {
			err = child.Render(context, writer)
			if err != nil {
				return false
			}
		}
		return true
	})
	return err
}
