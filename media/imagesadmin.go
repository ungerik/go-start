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
	return view.Views{
		UploadImageButton("", "", Config.Admin.ThumbnailSize, "function(){window.location.reload()}"),
		view.HR(),
		&view.Form{
			SubmitButtonText:  "Delete smaller image versions (they will be regenerated)",
			SubmitButtonClass: Config.Admin.ButtonClass,
			FormID:            "deleteVersions",
			OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
				var image Image

				for i := ImageIterator(); i.Next(&image); {
					err = image.DeleteVersions()
					if err != nil {
						return "", nil, err
					}
					err = image.Save()
					if err != nil {
						return "", nil, err
					}
				}
				return "", view.StringURL("."), nil
			},
		},
		&view.ModelIteratorView{
			GetModelIterator: func(ctx *view.Context) model.Iterator {
				return Config.Backend.ImageIterator()
			},
			GetModel: func(ctx *view.Context) (interface{}, error) {
				return new(Image), nil
			},
			GetModelView: func(ctx *view.Context, m interface{}) (view.View, error) {
				image := *m.(*Image) // copy by value because it will be used in a closure later on
				refCount, err := image.CountRefs()
				if err != nil {
					return nil, err
				}
				thumbnail, err := image.Thumbnail(Config.Admin.ThumbnailSize)
				if err != nil {
					return nil, err
				}
				deleteConfirmation := fmt.Sprintf("Are you sure you want to delete the image %s?", image.Title)
				if refCount > 0 {
					deleteConfirmation += fmt.Sprintf(" It is used %d times!", refCount)
				}
				editor := view.DIV(Config.Admin.ImageEditorClass,
					view.H3(image.Title.Get()),
					view.P(view.A_blank(image.FileURL(), "Link to original") /*, view.Printf(" | Used %d times", refCount)*/),
					view.P("Size in MB: ", view.Printf("%f", (float32)(image.Size.Get())/(1024.0*1024.0))),
					view.DIV(Config.Admin.ThumbnailFrameClass,
						thumbnail.View(""),
					),
					&view.Form{
						FormID:            "edit" + image.ID.Get(),
						SubmitButtonClass: Config.Admin.ButtonClass,
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
						SubmitButtonClass:   Config.Admin.ButtonClass,
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
