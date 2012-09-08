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

		ctx.Response.Printf(`{success: true, imageID: "%s", thumbnailURL: "%s", thumbnailSize: %d}`, i.ID, v.URL().URL(ctx), thumbnailSize)
		return nil
	},
))

// Uses https://github.com/valums/file-uploader
func UploadImageButton(thumbnailSize int, onComplete string) view.View {
	if onComplete == "" {
		onComplete = "function(){}"
	}

	const registerUploaderScript = `jQuery(function() {
		var uploader = new qq.FileUploader({
			debug: true,
		    element: jQuery("#%s")[0],
		    action: "%s",
		    allowedExtensions: ["png", "jpg", "jpeg", "gif", "bmp", "tif", "tiff"],
		    acceptFiles: ["image/png", "image/jpeg", "image/gif", "image/bmp", "image/tiff"],
		    multiple: false,
		    onComplete: %s
		});
	});`

	uploader := view.DIV("uploader")
	uploaderID := uploader.ID()
	uploaderScript := view.RenderView(
		func(ctx *view.Context) error {
			uploadURL := UploadImage.URL(ctx.ForURLArgsConvert(thumbnailSize))
			script := fmt.Sprintf(registerUploaderScript, uploaderID, uploadURL, onComplete)
			ctx.Response.RequireScript(script, 20)
			return nil
		},
	)
	return view.Views{
		&view.If{
			Condition: !Config.NoDynamicStyleAndScript,
			Content: view.Views{
				view.RequireStyleURL("/media/fileuploader.css", 0),
				view.RequireScriptURL("/media/fileuploader.js", 0),
				// view.RequireScriptURL("/media/media.js", 10),
			},
		},
		uploader,
		uploaderScript,
	}
}
