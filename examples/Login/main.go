package main

// import (
// 	"log"
// 	"os"

// 	// "github.com/ungerik/go-mail"
// 	"github.com/ungerik/go-start/errs"
// 	"github.com/ungerik/go-start/mongo"
// 	"github.com/ungerik/go-start/mongomedia"
// 	"github.com/ungerik/go-start/user"
// 	"github.com/ungerik/go-start/view"
// 	// "github.com/ungerik/go-start/debug"

// 	// "github.com/ungerik/go-start-start/models"
// 	"github.com/ungerik/go-start-start/views"

// 	// Dummy-import view packages for initialization:
// 	_ "github.com/ungerik/go-start-start/views/root"
// )

func main() {
	// ///////////////////////////////////////////////////////////////////////////
	// // Setup email
	// // todo replace with config file
	// // email.InitGmail("my@email.com", "password")

	// ///////////////////////////////////////////////////////////////////////////
	// // Config mongo
	// // todo replace with config file
	// err := mongo.InitLocalhost("gostart-example", "gostart-example", "")
	// errs.PanicOnError(err)
	// defer mongo.Close()

	// ///////////////////////////////////////////////////////////////////////////
	// // Config user
	// // todo replace with config file
	// user.Init(models.People)
	// doc, _, err := user.EnsureExists("admin", "my@email.com", "admin", true)
	// // errs.PanicOnError(err)
	// u := doc.(*models.Person)
	// u.SuperAdmin = true
	// err = u.Save()
	// errs.PanicOnError(err)

	// ///////////////////////////////////////////////////////////////////////////
	// // Config view
	// baseDirs := []string{
	// 	".",
	// 	"../../ungerik/go-start/",
	// }
	// err = view.Init("example.com", views.CookieSecret, baseDirs...)
	// errs.PanicOnError(err)
	// defer view.Close()

	// view.Config.Debug.Mode = true
	// view.Config.Debug.PrintPaths = false
	// view.Config.Debug.PrintRedirects = false
	// view.Config.DisableCachedViews = false
	// view.Config.RedirectSubdomains = []string{"www"}
	// view.Config.Page.DefaultMetaViewport = "width=1000px"
	// view.Config.Page.DefaultWriteScripts = view.PageWriters(
	// 	view.PageWritersFilterPort(80, view.GoogleAnalytics(views.GoogleAnalyticsID)),
	// )
	// view.Config.LoginSignupPage = &views.LoginSignup
	// //view.Config.GlobalAuth = view.NewBasicAuth("statuplive.in", "gostart", "gostart")
	// view.Config.OnPreAuth = func(context *view.Context) error {
	// 	user.OfSession(context) // Sets context.User
	// 	return nil
	// }

	// ///////////////////////////////////////////////////////////////////////////
	// // Run server
	// view.RunConfigFile(views.Paths(), "run.config")
}
