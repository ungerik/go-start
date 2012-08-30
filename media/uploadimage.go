package media

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ungerik/go-start/view"
)

var UploadImage = view.NewViewURLWrapper(view.RenderView(
	func(ctx *view.Context) error {
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
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return formatError(err)
		}
		i, err := NewImage(filename, b)
		if err != nil {
			return formatError(err)
		}

		ctx.Response.Printf(`{success: true, id: "%s"}`, i.ID)
		return nil
	},
))
