package main

import (
	"github.com/ungerik/go-start/config"
	"github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/ViewPaths/views"

	// Dummy-import all views sub-packages for initialization:
	_ "github.com/ungerik/go-start/examples/ViewPaths/views/admin"
	_ "github.com/ungerik/go-start/examples/ViewPaths/views/admin/user_0"
	_ "github.com/ungerik/go-start/examples/ViewPaths/views/root"
)

func main() {
	defer config.Close() // Close all packages on exit	
	config.Load("config.json", &view.Config)
	view.RunServer(views.Paths())
}
