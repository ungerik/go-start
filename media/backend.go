package media

import (
	"io"
)

type Backend interface {
	LoadImage(id string) (*Image, error)
	// image.ID will be updated if empty
	SaveImage(image *Image) error

	// ImageVersionReader returns an io.ReadCloser to read the image-data
	// with the given id from the backend
	// If there is no image with the given id,
	// err will be of type ErrInvalidImageID
	ImageVersionReader(id string) (reader io.ReadCloser, ctype string, err error)

	// ImageVersionWriter returns an io.WriteCloser to write the image-data
	// to the backend. version.ID can be empty for a new image or the id
	// of an existing image. version.ID can be changed by the function call
	// regardless of the former value
	ImageVersionWriter(version *ImageVersion) (writer io.WriteCloser, err error)
}

type ErrInvalidImageID string

func (self ErrInvalidImageID) Error() string {
	return "Invalid image ID: \"" + string(self) + "\""
}
