package media

import (
	"github.com/ungerik/go-start/view"
)

var AllThumbnailsAPI = view.NewViewURLWrapper(view.RenderViewBindURLArgs(
	func(ctx *view.Context, thumbnailSize int) error {
		// ctx.Response.SetContentType("application/json")
		ctx.Response.SetContentTypeByExt(".json")
		ctx.Response.WriteString("[\n")
		first := true
		i := Config.Backend.ImageIterator()
		for doc := i.Next(); doc != nil; doc = i.Next() {
			image := doc.(*Image)
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
				image.ID, image.TitleOrFilename(), thumbnail.GetURL().URL(ctx),
			)
		}
		if i.Err() != nil {
			return i.Err()
		}
		ctx.Response.WriteString("\n]")
		return nil
	},
))
