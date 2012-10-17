package media

import (
	"fmt"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type imageData struct {
	Title    model.String
	Filename model.String
	Link     model.Url
}

func ImagesAdmin() view.View {
	uploader := view.DIV("")
	uploaderID := uploader.ID()
	uploader.Content = UploadImageButton("#"+uploaderID, "", "", Config.ImagesAdmin.ThumbnailSize, "function(){window.location.reload()}")
	return view.Views{
		uploader,
		&view.ModelIteratorView{
			Model: new(Image),
			GetModelIterator: func(ctx *view.Context) model.Iterator {
				return Config.Backend.ImageIterator()
			},
			GetModelIteratorView: func(ctx *view.Context, m interface{}) (view.View, error) {
				image := m.(*Image)
				refCount, err := image.CountRefs()
				if err != nil {
					return nil, err
				}
				thumbnail, err := image.Thumbnail(Config.ImagesAdmin.ThumbnailSize)
				if err != nil {
					return nil, err
				}
				deleteConfirmation := fmt.Sprintf("Are you sure you want to delete the image %s?", image.TitleOrFilename())
				if refCount > 0 {
					deleteConfirmation += fmt.Sprintf(" It is used %d times!", refCount)
				}
				editor := view.DIV(Config.ImagesAdmin.ImageEditorClass,
					view.H3(image.TitleOrFilename()),
					view.P(
						view.A_blank(image.GetURL(), "Link to original"),
						view.Printf(" | Used %d times", refCount),
					),
					view.DIV(Config.ImagesAdmin.ThumbnailFrameClass,
						thumbnail.View(""),
					),
					&view.Form{
						FormID:            "edit" + image.ID.Get(),
						SubmitButtonClass: Config.ImagesAdmin.ButtonClass,
						GetModel: func(form *view.Form, ctx *view.Context) (interface{}, error) {
							return &imageData{
								Title:    image.Title,
								Filename: model.String(image.Filename()),
								Link:     image.Link,
							}, nil
						},
						OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
							model := formModel.(*imageData)
							image.Title = model.Title
							image.Link = model.Link
							for i := range image.Versions {
								image.Versions[i].Filename = model.Filename
							}
							err = image.Save()
							return "", view.StringURL("."), err
						},
					},
					&view.Form{
						SubmitButtonText:    "Delete",
						SubmitButtonConfirm: deleteConfirmation,
						SubmitButtonClass:   Config.ImagesAdmin.ButtonClass,
						FormID:              "delete" + image.ID.Get(),
						OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
							_, err = image.RemoveAllRefs()
							if err != nil {
								return "", nil, err
							}
							return "", view.StringURL("."), image.Delete()
						},
					},
					view.DivClearBoth(),
				)
				return editor, nil
			},
		},
	}
}
