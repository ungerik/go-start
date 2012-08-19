package views

import (
	. "github.com/ungerik/go-start/view"
	"github.com/ungerik/go-start/media"
)

var LoginSignup *Page
var ConfirmEmail *Page
var Logout ViewWithURL

var Homepage ViewWithURL

func Paths() *ViewPath {
	//basicAuth := NewBasicAuth("domain", "username", "password")
	return &ViewPath{View: Homepage, Sub: []ViewPath{
		media.ViewPath("media"),
		{Name: "style.css", View: CSS},
		{Name: "admin", View: Admin, Auth: Admin_Auth, Sub: []ViewPath{}},
		{Name: "login", View: LoginSignup, Sub: []ViewPath{
			{Name: "confirm", View: ConfirmEmail},
		}},
		{Name: "logout", View: Logout},
		{Name: "profile", View: Profile},
	}}
}
