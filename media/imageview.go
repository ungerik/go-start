package media

import (
	"io"

	"github.com/ungerik/go-start/view"
)

var ImageView = view.NewViewURLWrapper(view.RenderView(
	func(ctx *view.Context) error {
		reader, _, contentType, err := Config.Backend.FileReader(ctx.URLArgs[0])
		if err != nil {
			if _, ok := err.(ErrNotFound); ok {
				return view.NotFound(ctx.URLArgs[0] + "/" + ctx.URLArgs[1] + " not found")
			}
			return err
		}
		_, err = io.Copy(ctx.Response, reader)
		if err != nil {
			return err
		}
		err = reader.Close()
		if err != nil {
			return err
		}
		ctx.Response.Header().Set("Content-Type", contentType)
		return nil
	},
))
