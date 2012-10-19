package media

import (
	"io"

	"github.com/ungerik/go-start/model"
)

type Backend interface {
	LoadImage(id string) (*Image, error)
	TryLoadImage(id string) (*Image, bool, error)

	// SaveImage saves image and updates its ID if it is empty.
	SaveImage(image *Image) error

	DeleteImage(image *Image) error
	DeleteImageVersion(id string) error

	// ImageVersionReader returns an io.ReadCloser to read the image-data
	// with the given id from the backend.
	// If there is no image with the given id,
	// an error of type ErrInvalidImageID will be returned.
	ImageVersionReader(id string) (reader io.ReadCloser, ctype string, err error)

	// ImageVersionWriter returns an io.WriteCloser to write the image-data
	// to the backend. version.ID can be empty for a new image or the id
	// of an existing image. version.ID can be changed by the function call
	// regardless of the former value
	ImageVersionWriter(version *ImageVersion) (writer io.WriteCloser, err error)

	// ImageIterator iterates all images, use Next(*Image), not Next(**Image)
	ImageIterator() model.Iterator

	// CountImageRefs counts all ImageRef occurances with imageID
	// in all known databases.
	CountImageRefs(imageID string) (count int, err error)

	// RemoveAllImageRefs removes all ImageRef occurances with imageID
	// in all known databases.
	RemoveAllImageRefs(imageID string) (count int, err error)
}

type ErrInvalidImageID string

func (self ErrInvalidImageID) Error() string {
	return "Invalid image ID: \"" + string(self) + "\""
}
