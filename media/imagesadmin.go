package media

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

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
			return view.DIV(Config.ImagesAdmin.Class,
				view.H3(image.TitleOrFilename()),
				view.A_blank(image.URL(), image.URL().URL(ctx)),
				thumbnail.ViewImage(Config.ImagesAdmin.ThumbnailClass),
			), nil
		},
	}
}
