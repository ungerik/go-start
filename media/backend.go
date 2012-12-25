package media

import (
	"io"

	"github.com/ungerik/go-start/model"
)

type Backend interface {
	// General file methods:

	FileWriter(filename, contentType string) (writer io.WriteCloser, id string, err error)
	// Returns ErrNotFound if no file with id is found.
	FileReader(id string) (reader io.ReadCloser, filename, contentType string, err error)
	// Returns ErrNotFound if no file with id is found.
	DeleteFile(id string) error

	// Blob methods:

	LoadBlob(id string) (*Blob, error)
	SaveBlob(blob *Blob) error
	// DeleteBlob does not delete the file associated with it, also use DeleteFile().
	DeleteBlob(blob *Blob) error

	// BlobIterator returns an iterator that iterates
	// all blobs as Blob structs.
	BlobIterator() model.Iterator

	// CountBlobRefs counts all BlobRef occurrences with blobID
	// in all known databases.
	CountBlobRefs(blobID string) (count int, err error)

	// RemoveAllBlobRefs removes all BlobRef occurrences with blobID
	// in all known databases.
	RemoveAllBlobRefs(blobID string) (count int, err error)

	// Image methods:

	// Returns ErrNotFound if no image with id is found.
	LoadImage(id string) (*Image, error)

	// SaveImage saves image and updates its ID if it is empty.
	SaveImage(image *Image) error

	DeleteImage(image *Image) error

	// ImageIterator returns an iterator that iterates
	// all images as Image structs.
	ImageIterator() model.Iterator

	// CountImageRefs counts all ImageRef occurrences with imageID
	// in all known databases.
	CountImageRefs(imageID string) (count int, err error)

	// RemoveAllImageRefs removes all ImageRef occurrences with imageID
	// in all known databases.
	RemoveAllImageRefs(imageID string) (count int, err error)
}

type ErrNotFound string

func (self ErrNotFound) Error() string {
	return "Media not found: \"" + string(self) + "\""
}
