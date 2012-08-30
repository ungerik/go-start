package media

import (
	"io"

	"github.com/ungerik/go-start/view"
)

var ContentView = view.NewViewURLWrapper(view.RenderView(
	func(response *view.Response) error {
		reader, contentType, err := Config.Backend.ImageVersionReader(response.Request.URLArgs[0])
		if err != nil {
			if _, ok := err.(ErrInvalidImageID); ok {
				return view.NotFound(response.Request.URLArgs[0] + "/" + response.Request.URLArgs[1] + " not found")
			}
			return err
		}
		_, err = io.Copy(response, reader)
		if err != nil {
			return err
		}
		err = reader.Close()
		if err != nil {
			return err
		}
		response.Header().Set("Content-Type", contentType)
		return nil
	},
))

func ContentViewPath(name string) view.ViewPath {
	return view.ViewPath{Name: name, Args: 2, View: ContentView}
}
