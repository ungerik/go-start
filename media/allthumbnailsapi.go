package media

import (
	"github.com/ungerik/go-start/view"
)

var AllThumbnailsAPI = view.NewViewURLWrapper(view.RenderViewBindURLArgs(
	func(ctx *view.Context, thumbnailSize int) error {
		ctx.Response.Header().Set("Content-Type", "application/json")
		ctx.Response.WriteString("[\n")
		first := true
		i := Config.Backend.ImageIterator()
		var image Image
		for i.Next(&image) {
			thumbnail, err := image.Thumbnail(thumbnailSize)
			if err != nil {
				return err
			}
			if first {
				first = false
			} else {
				ctx.Response.WriteString(",\n")
			}
			ctx.Response.Printf(
				`{"id": "%s", "title": "%s", "url": "%s"}`,
				image.ID, image.TitleOrFilename(), thumbnail.FileURL().URL(ctx),
			)
		}
		if i.Err() != nil {
			return i.Err()
		}
		ctx.Response.WriteString("\n]")
		return nil
	},
))
