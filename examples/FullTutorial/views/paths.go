package views

import (
	"github.com/ungerik/go-start/media"
	// "github.com/ungerik/go-start/debug"
	. "github.com/ungerik/go-start/view"
)

var (
	CSS                View
	Homepage           ViewWithURL
	LoginSignup        *Page
	Logout             ViewWithURL
	ConfirmEmail       *Page
	Profile            *Page
	Admin              *Page
	Admin_Auth         Authenticator
	Admin_Users        *Page
	Admin_ExportEmails ViewWithURL
	Admin_UserX        *Page
	Admin_Images       *Page
)

func Paths() *ViewPath {
	return &ViewPath{View: Homepage, Sub: []ViewPath{
		media.ViewPath("media"),
		{Name: "style.css", View: CSS},
		{Name: "admin", View: Admin, Auth: Admin_Auth, Sub: []ViewPath{
			{Name: "users", View: Admin_Users, Auth: Admin_Auth, Sub: []ViewPath{
				{Name: "export-emails", View: Admin_ExportEmails, Auth: Admin_Auth},
			}},
			{Name: "user", Args: 1, View: Admin_UserX, Auth: Admin_Auth},
			{Name: "images", View: Admin_Images, Auth: Admin_Auth},
		}},
		{Name: "login", View: LoginSignup, Sub: []ViewPath{
			{Name: "confirm", View: ConfirmEmail},
		}},
		{Name: "logout", View: Logout},
		{Name: "profile", View: Profile},
	}}
}
