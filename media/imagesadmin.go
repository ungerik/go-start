package media

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type imageData struct {
	Title    model.String
	Filename model.String
	Link     model.Url
}

func ImagesAdmin() *view.ModelIteratorView {
	return &view.ModelIteratorView{
		GetModelIterator: func(ctx *view.Context) model.Iterator {
			return Config.Backend.ImageIterator()
		},
		GetModelIteratorView: func(ctx *view.Context, m interface{}) (view.View, error) {
			image := m.(*Image)
			thumbnail, err := image.Thumbnail(Config.ImagesAdmin.ThumbnailSize)
			if err != nil {
				return nil, err
			}
			v := view.DIV(Config.ImagesAdmin.Class,
				view.H3(image.TitleOrFilename()),
				thumbnail.ViewImage(Config.ImagesAdmin.ThumbnailClass),
				view.HTML("Image URL: "),
				view.A_blank(image.URL(), image.URL().URL(ctx)),
				&view.Form{
					FormID:            "image" + image.ID.Get(),
					SubmitButtonClass: Config.ImagesAdmin.ButtonClass,
					GetModel: func(form *view.Form, ctx *view.Context) (model interface{}, err error) {
						return &imageData{
							Title:    image.Title,
							Filename: image.Versions[0].Filename,
							Link:     image.Link,
						}, nil
					},
					OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
						model := formModel.(*imageData)
						image.Title = model.Title
						image.Link = model.Link
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
			)
			return v, nil
		},
	}
}
