package media

import (
	"io"
	"github.com/ungerik/go-start/view"
	"github.com/ungerik/go-start/utils"
)

var View view.ViewWithURL = view.NewViewURLWrapper(view.RenderView(
	func(context *view.Context, writer *utils.XMLWriter) error {
		reader, contentType, err := Config.Backend.ImageVersionReader(context.PathArgs[0])
		if err != nil {
			if _, ok := err.(ErrInvalidImageID); ok {
				return view.NotFound(context.PathArgs[0] + "/" + context.PathArgs[1] + " not found")
			}
			return err
		}
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
		err = reader.Close()
		if err != nil {
			return err
		}
		context.Header().Set("Content-Type", contentType)
		return nil
	},
))

func ViewPath(name string) view.ViewPath {
	return view.ViewPath{Name: name, Args: 2, View: View}
}
