package admin

import (
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Admin_ExportEmails = NewViewURLWrapper(
		RenderView(
			func(ctx *Context) (err error) {
				// Download instead of display:
				// ctx.Response.ContentDispositionAttachment("emails.csv")
				i := models.Users.Iterator()
				var u models.User
				for i.Next(&u) {
					if len(u.Email) > 0 {
						ctx.Response.Printf("%s, %s, %s \n", u.Email[0].Address, u.Name.First.String(), u.Name.Last.String())
					}
				}
				return i.Err()
			},
		),
	)
}
