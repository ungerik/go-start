package admin

import (
	. "github.com/ungerik/go-start/view"

	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Admin = NewAdminPage("Admin Dashboard",
		&Form{
			FormID:            "clearcaches",
			SubmitButtonText:  "Clear Page Caches",
			SubmitButtonClass: "button",
			OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
				ClearAllCaches()
				return "All caches cleared", nil, nil
			},
		},
	)
}
