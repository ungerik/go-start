package media

import (
	"github.com/ungerik/go-start/view"
)

var AllThumbnailsView = view.NewViewURLWrapper(view.RenderViewBindURLArgs(
	func(ctx *view.Context, thumbnailSize int) error {
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
			ctx.Response.Printf(`{id: "%s", url: "%s"}`, image.ID, thumbnail.GetURL().URL(ctx))
		}
		if i.Err() != nil {
			return i.Err()
		}
		ctx.Response.WriteString("\n]")
		return nil
	},
))
