package root

import (
	"github.com/ungerik/go-start/user"
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Profile = NewPublicPage("My Profile | go-start Tutorial",
		DIV("public-content",
			DynamicView(
				func(ctx *Context) (view View, err error) {
					var usr models.User
					found, err := user.OfSession(ctx.Session, &usr)
					if err != nil {
						return nil, err
					}
					if !found {
						return H1("You have to be logged in to edit your startup"), nil
					}
					email := usr.PrimaryEmail()

					view = DIV("row",
						DIV("cell right-border",
							TitleBar("My Profile"),
							DIV("main-content",
								H3("Email: ", email),
								H3("Name:"),
								P(&Form{
									SubmitButtonText:  "Save name",
									SubmitButtonClass: "button",
									FormID:            "profile",
									GetModel: func(form *Form, ctx *Context) (interface{}, error) {
										return &usr.Name, nil
									},
									OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
										return "", StringURL("."), usr.Save()
									},
								}),
								H3("Password:"),
								P(&Form{
									SubmitButtonText:  "Save password",
									SubmitButtonClass: "button",
									FormID:            "password",
									GetModel: func(form *Form, ctx *Context) (interface{}, error) {
										return new(user.PasswordFormModel), nil
									},
									OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
										m := formModel.(*user.PasswordFormModel)
										usr.Password.SetHashed(m.Password1.Get())
										return "", StringURL("."), usr.Save()
									},
								}),
							),
						),
						DivClearBoth(),
					)
					return view, nil
				},
			),
		),
	)
}
