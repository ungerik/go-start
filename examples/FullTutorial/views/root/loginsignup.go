package root

import (
	// "github.com/ungerik/go-start/model"
	gostartuser "github.com/ungerik/go-start/user"
	. "github.com/ungerik/go-start/view"

	// "github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Logout = NewViewURLWrapper(gostartuser.LogoutView(nil))

	ConfirmEmail = NewPublicPage("Email Confirmation | go-start Tutorial",
		DIV("public-content",
			DIV("main",
				TitleBar("Email confirmation"),
				DIV("main-content", gostartuser.EmailConfirmationView(IndirectURL(&Profile))),
			),
		),
	)

	LoginSignup = NewPublicPage("Login or Sign up | go-start Tutorial",
		DIV("public-content",
			DynamicView(
				func(ctx *Context) (view View, err error) {
					_, hasFrom := ctx.Request.Params["from"]
					id := ctx.Session.ID()
					if id == "" && hasFrom {
						view = DIV("main",
							DIV("main-content",
								H3("Your account doesn't have sufficient rights to view this page"),
								Printf("You may <a href='%s'>logout</a> and login with a different account", Logout.URL(ctx)),
							),
						)
					} else {
						view = DIV("row",
							DIV("cell left",
								TitleBar("Log in"),
								DIV("main-content",
									gostartuser.NewLoginForm("Log in", "login", "error", "success", IndirectURL(&Homepage)),
								),
							),
							DIV("cell right",
								TitleBarRight("Sign up"),
								DIV("main-content",
									gostartuser.NewSignupForm("Sign up", "signup", "error", "success", IndirectURL(&ConfirmEmail), nil),
								),
							),
							DivClearBoth(),
						)
					}
					return view, nil
				},
			),
		),
	)
}
