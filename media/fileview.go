package media

import (
	"io"

	"github.com/ungerik/go-start/config"
	"github.com/ungerik/go-start/view"
)

var FileView = view.NewViewURLWrapper(view.RenderView(
	func(ctx *view.Context) error {
		reader, _, contentType, err := Config.Backend.FileReader(ctx.URLArgs[0])
		if err != nil {
			if _, ok := err.(ErrNotFound); ok {
				err = view.NotFound(ctx.URLArgs[0] + "/" + ctx.URLArgs[1] + " not found")
				config.Logger.Println("FileView:", err)
				return err
			}
			return err
		}
		defer reader.Close()
		_, err = io.Copy(ctx.Response, reader)
		if err != nil {
			config.Logger.Println("FileView:", err)
			return err
		}
		ctx.Response.Header().Set("Content-Type", contentType)
		return nil
	},
))
