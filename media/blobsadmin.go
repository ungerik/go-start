package media

import (
	"fmt"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type blobData struct {
	Title    model.String
	Filename model.String
	Link     model.Url
}

func BlobsAdmin() view.View {
	return view.Views{
		UploadBlobButton("", "", "function(){window.location.reload()}"),
		&view.ModelIteratorView{
			GetModelIterator: func(ctx *view.Context) model.Iterator {
				// debug.Dump(Config.Backend.Blobs.Count())
				return Config.Backend.BlobIterator()
			},
			GetModel: func(ctx *view.Context) (interface{}, error) {
				return new(Blob), nil
			},
			GetModelView: func(ctx *view.Context, m interface{}) (view.View, error) {
				blob := m.(*Blob)
				// refCount, err := blob.CountRefs()
				// if err != nil {
				// 	return nil, err
				// }
				deleteConfirmation := fmt.Sprintf("Are you sure you want to delete the blob %s?", blob.TitleOrFilename())
				// if refCount > 0 {
				// 	deleteConfirmation += fmt.Sprintf(" It is used %d times!", refCount)
				// }
				editor := view.DIV(Config.Admin.ImageEditorClass,
					view.H3(blob.TitleOrFilename()),
					view.P(
						view.A_blank(blob.FileURL(), "Link to file"),
						// view.Printf(" | Used %d times", refCount),
					),
					&view.Form{
						FormID:            "edit" + blob.ID.Get(),
						SubmitButtonClass: Config.Admin.ButtonClass,
						GetModel: func(form *view.Form, ctx *view.Context) (interface{}, error) {
							return &blobData{
								Title:    blob.Title,
								Filename: blob.Filename,
								Link:     blob.Link,
							}, nil
						},
						OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
							model := formModel.(*blobData)
							blob.Title = model.Title
							blob.Filename = model.Filename
							blob.Link = model.Link
							err = blob.Save()
							return "", view.StringURL("."), err
						},
					},
					&view.Form{
						SubmitButtonText:    "Delete",
						SubmitButtonConfirm: deleteConfirmation,
						SubmitButtonClass:   Config.Admin.ButtonClass,
						FormID:              "delete" + blob.ID.Get(),
						OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
							// _, err = blob.RemoveAllRefs()
							// if err != nil {
							// 	return "", nil, err
							// }
							return "", view.StringURL("."), blob.Delete()
						},
					},
					view.DivClearBoth(),
				)
				return editor, nil
			},
		},
	}
}
