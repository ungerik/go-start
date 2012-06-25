package media

import (
	"io"
)

type BackendImage struct {
	Reader      io.ReadCloser
	Filename    string
	ContentType string
	Width       int
	Height      int
	Grayscale   bool
}

type Backend interface {
	Image(id string) (*Image, error)

	// The image data can be read from image.Reader
	LoadImage(id string) (image *BackendImage, found bool, err error)
	// The image data will be copied from image.Reader
	SaveImage(id string, image *BackendImage) (err error)
	// The image data will be copied from image.Reader
	SaveNewImage(image *BackendImage) (id string, err error)
}

// type ErrInvalidImageID string

// func (self ErrInvalidImageID) Error() string {
// 	return "Invalid image ID: \"" + string(self) + "\""
// }
