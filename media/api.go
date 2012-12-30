package media

import (
	// "strings"

	// "github.com/ungerik/go-start/utils"
	"github.com/ungerik/go-start/view"
)

type Api struct {
	AllThumbnails view.ViewWithURL
	// AllBlobs      view.ViewWithURL
}

var API = Api{

	AllThumbnails: view.RenderViewWithURLBindURLArgs(
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
					image.ID, image.Title, thumbnail.FileURL().URL(ctx),
				)
			}
			if i.Err() != nil {
				return i.Err()
			}
			ctx.Response.WriteString("\n]")
			return nil
		},
	),

	// AllBlobs: view.RenderViewWithURL(
	// 	func(ctx *view.Context) error {
	// 		if !Config.NoDynamicStyleAndScript {
	// 			ctx.Response.RequireScript(string(view.JQuery), 0)
	// 			ctx.Response.RequireScript(string(view.JQueryUI), 1)
	// 		}
	// 		searchTerm, _ := ctx.Request.Params["term"]
	// 		searchTerm = strings.ToLower(searchTerm)

	// 		ctx.Response.Header().Set("Content-Type", "application/json")
	// 		ctx.Response.WriteByte('[')
	// 		first := true
	// 		i := BlobIterator()
	// 		var blob Blob
	// 		for i.Next(&blob) {
	// 			if first {
	// 				first = false
	// 			} else {
	// 				ctx.Response.WriteByte(',')
	// 			}
	// 			ctx.Response.Printf(`{"id": "%s", "title": "%s"}`, blob.ID, blob.Title)
	// 		}
	// 		ctx.Response.WriteByte(']')
	// 		return i.Err()
	// 	},
	// ),
}
