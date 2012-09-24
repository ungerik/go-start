package views

import (
	. "github.com/ungerik/go-start/view"
)

/*
Declare variables for all views.
Will be initialized by sub-packages.
Declared here to make them accessible from all sub-packages,
for instance to get their URL.
*/
var (
	Homepage ViewWithURL

	Admin       *Page
	Admin_User0 *Page // User0 means URL argument 0 defines the user for the view
	AdminAuth   Authenticator

	GetXML  ViewWithURL
	GetJSON ViewWithURL
)

// Paths() returns the URL structure of the site
func Paths() *ViewPath {
	return &ViewPath{View: Homepage, Sub: []ViewPath{
		{Name: "get.xml", View: GetXML},
		{Name: "get.json", View: GetJSON},
		// Paths can have an Authenticator
		{Name: "admin", View: Admin, Auth: AdminAuth, Sub: []ViewPath{
			// Args is the number arguments that get parsed from the URL.
			// URL for empty Name: /admin/<URLArg0>/
			// URL for Name = "user": /admin/user/<URLArg0>/
			{Args: 1, View: Admin_User0, Auth: AdminAuth},
		}},
	}}
}
