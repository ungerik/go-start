package media

import (
	"github.com/ungerik/go-start/view"
	"github.com/ungerik/go-start/utils"
	"io"
	// "path"
)

var View view.ViewWithURL = view.NewViewURLWrapper(view.RenderView(
	func(context *view.Context, writer *utils.XMLWriter) error {
		image, found, err := Config.Backend.LoadImage(context.PathArgs[0])
		if err != nil {
			return err
		}
		if !found {
			return view.NotFound(context.PathArgs[0] + "/" + context.PathArgs[1] + " not found")
		}
		context.Header().Set("Content-Type", image.ContentType)
		defer image.Reader.Close()
		_, err = io.Copy(writer, image.Reader)
		return err
	},
))

func ViewPath(name string) view.ViewPath {
	return view.ViewPath{Name: "media", Args: 2, View: View}
}
