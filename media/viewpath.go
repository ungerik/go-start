package media

import (
	"github.com/ungerik/go-start/view"
)

// ViewPath returns the view.ViewPath for all media URLs.
func ViewPath(name string) view.ViewPath {
	return view.ViewPath{Name: name, Sub: []view.ViewPath{
		{Name: "file", Args: 2, View: FileView},
		{Name: "upload-blob", View: UploadBlob},
		{Name: "upload-image", Args: 1, View: UploadImage},
		{Name: "thumbnails.json", Args: 1, View: API.AllThumbnails},
		// {Name: "blobs.json", View: API.AllBlobs},
	}}
}
