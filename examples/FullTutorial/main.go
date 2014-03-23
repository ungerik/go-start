package main

import (
	"github.com/ungerik/go-mail"
	"github.com/ungerik/go-start/config"
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/mongomedia"
	"github.com/ungerik/go-start/user"
	"github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/FullTutorial/models"
	"github.com/ungerik/go-start/examples/FullTutorial/views"

	// Dummy-import view packages for initialization:
	_ "github.com/ungerik/go-start/examples/FullTutorial/views/admin"
	_ "github.com/ungerik/go-start/examples/FullTutorial/views/root"
)

func main() {
	debug.Nop()

	///////////////////////////////////////////////////////////////////////////
	// Load configuration

	defer config.Close() // Close all packages on exit

	config.Load("config.json",
		&EmailConfig{&email.Config},
		&mongo.Config,
		&user.Config,
		&view.Config,
		&media.Config,
		&mongomedia.Config,
	)

	///////////////////////////////////////////////////////////////////////////
	// Ensure that an admin user exists

	var admin models.User
	_, err := user.EnsureExists("admin", "mail@ungerik.net", "test", true, &admin)
	errs.PanicOnError(err)
	admin.Admin = true
	err = admin.Save()
	errs.PanicOnError(err)

	///////////////////////////////////////////////////////////////////////////
	// Config view

	view.Config.NamedAuthenticators["admin"] = views.Admin_Auth
	view.Config.Debug.Mode = true
	// view.Config.Debug.LogPaths = true
	// view.Config.Debug.LogRedirects = true
	// view.Config.DisableCachedViews = true
	// view.Config.GlobalAuth = view.NewBasicAuth("statuplive.in", "gostart", "gostart")

	///////////////////////////////////////////////////////////////////////////
	// Run server

	view.RunServer(views.Paths())
}

///////////////////////////////////////////////////////////////////////////
// Extend email.Configuration with a Name() method
// to make it compatible with config.Load()

type EmailConfig struct {
	*email.Configuration
}

func (self *EmailConfig) Name() string {
	return "email"
}
