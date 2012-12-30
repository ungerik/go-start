package media

import (
	"fmt"
	"io"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/config"
	"github.com/ungerik/go-start/view"
)

var UploadBlob = view.NewViewURLWrapper(view.RenderViewBindURLArgs(
	func(ctx *view.Context, thumbnailSize int) error {
		formatError := func(err error) error {
			config.Logger.Println("UploadBlob:", err)
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
		blob, err := NewBlobFromReader(filename, r)
		if err != nil {
			return formatError(err)
		}

		ctx.Response.Printf(`{success: true, blobID: "%s", blobSize: %d}`, blob.ID, blob.Size)
		return nil
	},
))

// todo move to media.js
func uploadBlobButtonScript(ctx *view.Context, parentSelector, dropZoneSelector, listSelector, onComplete string) string {
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
	return fmt.Sprintf(
		`jQuery(function() {
			var uploader = new qq.FileUploader({
				debug: true,
			    element: jQuery("%s")[0],
			    extraDropzones: %s,
			    listElement: %s,
			    action: "%s",
			    sizeLimit: 1024*1024*64,
			    multiple: false,
			    onComplete: %s
			});
		});`,
		parentSelector,
		extraDropzones,
		listElement,
		UploadBlob.URL(ctx),
		onComplete,
	)
}

// Uses https://github.com/valums/file-uploader
func RequireUploadBlobButtonScript(parentSelector, dropZoneSelector, listSelector, onComplete string) view.View {
	return view.RenderView(
		func(ctx *view.Context) error {
			script := uploadBlobButtonScript(ctx, parentSelector, dropZoneSelector, listSelector, onComplete)
			ctx.Response.RequireScript(script, 20)
			if !Config.NoDynamicStyleAndScript {
				ctx.Response.RequireStyleURL("/media/fileuploader.css", 0)
				ctx.Response.RequireScriptURL("/media/fileuploader.js", 0)
			}
			return nil
		},
	)
}

func UploadBlobButton(dropZoneSelector, listSelector, onComplete string) view.View {
	var button view.Div
	button.Content = view.Views{
		view.H1("jQuery required!"),
		RequireUploadBlobButtonScript("#"+button.ID(), dropZoneSelector, listSelector, onComplete),
	}
	return &button
}
