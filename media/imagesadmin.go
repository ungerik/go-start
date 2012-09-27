package media

import (
	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type imageData struct {
	Title    model.String
	Filename model.String
	Link     model.Url
}

func ImagesAdmin() view.View {
	return view.Views{
		UploadImageButton(Config.ImagesAdmin.ThumbnailSize, "function(){window.location.reload()}"),
		&view.ModelIteratorView{
			GetModelIterator: func(ctx *view.Context) model.Iterator {
				debug.Dump(Config.Backend)
				return Config.Backend.ImageIterator()
			},
			GetModelIteratorView: func(ctx *view.Context, m interface{}) (view.View, error) {
				image := m.(*Image)
				thumbnail, err := image.Thumbnail(Config.ImagesAdmin.ThumbnailSize)
				if err != nil {
					return nil, err
				}
				editor := view.DIV(Config.ImagesAdmin.ImageEditorClass,
					view.H3(image.TitleOrFilename()),
					view.A_blank(image.GetURL(), "Link to original"),
					view.DIV(Config.ImagesAdmin.ThumbnailFrameClass,
						thumbnail.ViewImage(""),
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
						SubmitButtonConfirm: "Are you sure you want to delete the image " + image.TitleOrFilename() + "?",
						SubmitButtonClass:   Config.ImagesAdmin.ButtonClass,
						FormID:              "delete" + image.ID.Get(),
						OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
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
