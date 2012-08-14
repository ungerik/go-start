package view

import (
	"io"
	"log"
	"strconv"

	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/templatesystem"
	"github.com/ungerik/go-start/utils"
	// "github.com/ungerik/goconfig"
	"github.com/ungerik/web.go"

//	"github.com/ungerik/go-start/debug"
)

var viewIdChan = make(chan string, 16)
var viewIdCounter int64

func init() {
	go func() {
		for {
			// Use base32 encoding for ids to make them shorter
			viewIdChan <- strconv.FormatInt(viewIdCounter, 32)
			viewIdCounter++
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

func RunServer(paths *ViewPath) {
	addr := Config.ListenAndServeAt
	if !Config.IsProductionServer && Config.Debug.ListenAndServeAt != "" {
		addr = Config.Debug.ListenAndServeAt
	}
	log.Print("IsProductionServer = ", Config.IsProductionServer)
	log.Print("Debug.Mode = ", Config.Debug.Mode)
	paths.initAndRegisterViewsRecursive("/")
	web.Run(addr)
}

// func RunConfigFile(paths *ViewPath, filename string) {
// 	config, err := config.ReadDefault(filename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var host string
// 	host, err = config.String("server", "host")
// 	if err != nil {
// 		panic(err)
// 	}

// 	var port string
// 	port, err = config.String("server", "port")
// 	if err != nil {
// 		panic(err)
// 	}

// 	var baseURL string
// 	baseURL, err = config.String("server", "url")
// 	if err != nil {
// 		panic(err)
// 	}

// 	var recoverPanic bool
// 	recoverPanic, _ = config.Bool("server", "recoverPanic")

// 	Run(paths, host+":"+port, baseURL, recoverPanic)
// }

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

func RenderChildViewsHTML(parent View, response *Response) (err error) {
	parent.IterateChildren(func(parent View, child View) (next bool) {
		if child != nil {
			err = child.Render(response)
			if err != nil {
				return false
			}
		}
		return true
	})
	return err
}
