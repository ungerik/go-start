package admin

import (
	"github.com/ungerik/go-start/config"
	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/utils"
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Admin_Users = NewAdminPage("Admin Users",
		Views{
			A(Admin_ExportEmails, "Download all emails as .csv"),
			HR(),
			&Form{
				SubmitButtonText:  "Add Person",
				SubmitButtonClass: "button",
				FormID:            "addperson",
				OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
					var u models.User
					models.Users.InitDocument(&u)
					u.Name.First = "[First]"
					u.Name.Last = "[Last]"
					u.AddEmail("", "")
					err := u.Save()
					if err != nil {
						return "", nil, err
					}
					return "", NewURLWithArgs(Admin_UserX, u.ID.Hex()), nil
				},
			},
			HR(),
			&ModelIteratorTableView{
				Class: "visual-table",
				GetModelIterator: func(ctx *Context) model.Iterator {
					// return models.Users.Iterator()
					return models.Users.SortFunc(
						func(a, b *models.User) bool {
							return utils.CompareCaseInsensitive(a.Name.String(), b.Name.String())
						},
					)
				},
				GetRowModel: func(ctx *Context) (interface{}, error) {
					return new(models.User), nil
				},
				GetHeaderRowViews: func(ctx *Context) (views Views, err error) {
					return Views{HTML("Nr"), HTML("Name"), HTML("Email"), HTML("Edit"), HTML("Delete")}, nil
				},
				GetRowViews: func(row int, rowModel interface{}, ctx *Context) (views Views, err error) {
					u := rowModel.(*models.User)
					editURL := Admin_UserX.URL(ctx.ForURLArgs(u.ID.Hex()))
					return Views{
						Printf("%d", row+1),
						Escape(u.Name.String()),
						Escape(u.PrimaryEmail()),
						A(editURL, "Edit"),
						&Form{
							SubmitButtonText:    "Delete",
							SubmitButtonConfirm: "Are you sure you want to delete " + u.Name.String() + "?",
							FormID:              "delete" + u.ID.Hex(),
							OnSubmit: func(form *Form, formModel interface{}, ctx *Context) (string, URL, error) {
								config.Logger.Printf("Deleting user '%s' with ID %s", u.Name.String(), u.ID.Hex())
								config.Logger.Printf("FormID: " + form.FormID)
								return "", StringURL("."), u.Delete()
							},
						},
					}, nil
				},
			},
		},
	)
}
