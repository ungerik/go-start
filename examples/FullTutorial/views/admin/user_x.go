package admin

import (
	"labix.org/v2/mgo/bson"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	// "github.com/ungerik/go-start/user"
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Admin_UserX = NewAdminPage("Admin Users",
		DynamicView(
			func(ctx *Context) (view View, err error) {
				var u models.User
				id := bson.ObjectIdHex(ctx.URLArgs[0])
				found, err := models.Users.TryDocumentWithID(id, &u)
				if err != nil {
					return nil, err
				}
				if !found {
					return nil, NotFound("404: User not found")
				}

				views := Views{
					H2(u.Name.String()),
					Printf("Email confirmed: %v", u.EmailPasswordConfirmed()),
					HR(),
					&Form{
						SubmitButtonText:  "Save password and mark email as confirmed",
						SubmitButtonClass: "button",
						FormID:            "password",
						DisabledFields:    []string{"Current_password"},
						GetModel: func(form *Form, ctx *Context) (interface{}, error) {
							return &passwordFormModel{Current_password: u.Password}, nil
						},
						OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
							m := formModel.(*passwordFormModel)
							u.Password.SetHashed(m.New_password.Get())
							u.ConfirmEmailPassword()
							return "", StringURL("."), u.Save()
						},
					},
					HR(),
					&Form{
						SubmitButtonText:         "Save User Data",
						SubmitButtonClass:        "button",
						FormID:                   "user" + u.ID.Hex(),
						GetModel:                 FormModel(&u),
						GeneralErrorOnFieldError: true,
						OnSubmit:                 OnFormSubmitSaveModelAndRedirect(Admin_Users),
					},
				}
				return views, nil
			},
		),
	)
}

type passwordFormModel struct {
	Current_password model.Password
	New_password     model.String `model:"minlen=6"`
	//Send_notification_email model.Bool
}
