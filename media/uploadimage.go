package media

import (
	"fmt"
	"io"
	"io/ioutil"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/view"
)

var UploadImage = view.NewViewURLWrapper(view.RenderViewBindURLArgs(
	func(ctx *view.Context, thumbnailSize int) error {
		formatError := func(err error) error {
			return fmt.Errorf(`{success: false, error: "%s"}`, err.Error())
		}

		filename := ctx.Request.Header.Get("X-File-Name")
		var r io.ReadCloser
		if filename != "" {
			r = ctx.Request.Body
		} else {
			f, h, err := ctx.Request.FormFile("qqfile")
			if err != nil {
				return formatError(err)
			}
			filename = h.Filename
			r = f
		}
		defer r.Close()
		// debug.Print("UploadImage", filename, thumbnailSize)
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return formatError(err)
		}
		i, err := NewImage(filename, b)
		if err != nil {
			return formatError(err)
		}
		v, err := i.Thumbnail(thumbnailSize)
		if err != nil {
			return formatError(err)
		}

		ctx.Response.Printf(`{success: true, imageID: "%s", thumbnailURL: "%s", thumbnailSize: %d}`, i.ID, v.GetURL().URL(ctx), thumbnailSize)
		return nil
	},
))

// Uses https://github.com/valums/file-uploader
func UploadImageButton(parentSelector, dropZoneSelector, listSelector string, thumbnailSize int, onComplete string) view.View {
	extraDropzones := "[]"
	if dropZoneSelector != "" {
		dropZoneSelector = fmt.Sprintf(`[jQuery("%s")[0]]`, dropZoneSelector)
	}
	listElement := "null"
	if listSelector != "" {
		listElement = fmt.Sprintf(`jQuery("%s")[0]`, listSelector)
	}
	if onComplete == "" {
		onComplete = "function(){}"
	}

	uploaderScript := view.RenderView(
		func(ctx *view.Context) error {
			uploadURL := UploadImage.URL(ctx.ForURLArgsConvert(thumbnailSize))
			script := fmt.Sprintf(
				`jQuery(function() {
					var uploader = new qq.FileUploader({
						debug: true,
					    element: jQuery("%s")[0],
					    extraDropzones: %s,
					    listElement: %s,
					    action: "%s",
					    allowedExtensions: ["png", "jpg", "jpeg", "gif", "bmp", "tif", "tiff"],
					    acceptFiles: ["image/png", "image/jpeg", "image/gif", "image/bmp", "image/tiff"],
					    sizeLimit: 1024*1024*64,
					    multiple: false,
					    onComplete: %s
					});
				});`,
				parentSelector,
				extraDropzones,
				listElement,
				uploadURL,
				onComplete,
			)
			ctx.Response.RequireScript(script, 20)
			if !Config.NoDynamicStyleAndScript {
				ctx.Response.RequireStyleURL("/media/fileuploader.css", 0)
				ctx.Response.RequireScriptURL("/media/fileuploader.js", 0)
			}
			return nil
		},
	)
	return uploaderScript
}
