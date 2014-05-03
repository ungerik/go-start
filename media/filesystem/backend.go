package filesystem

import (
	"io"
	"mime"
	"os"
	"path"

	"github.com/ungerik/go-start/media"
	"github.com/ungerik/go-start/model"
)

type Backend struct {
	BaseDir string
}

func (backend *Backend) FileWriter(filename, contentType string) (writer io.WriteCloser, id string, err error) {
	writer, err = os.Create(path.Join(backend.BaseDir, filename))
	return writer, filename, err
}

// Returns ErrNotFound if no file with id is found.
func (backend *Backend) FileReader(id string) (reader io.ReadCloser, filename, contentType string, err error) {
	reader, err = os.OpenFile(path.Join(backend.BaseDir, filename), os.O_RDONLY, 0660)
	return reader, id, mime.TypeByExtension(id), err
}

// Returns ErrNotFound if no file with id is found.
func (backend *Backend) DeleteFile(id string) error {
	return os.Remove(path.Join(backend.BaseDir, id))
}

func (backend *Backend) LoadBlob(id string) (*media.Blob, error) {
	panic("not implemented")
}

func (backend *Backend) SaveBlob(blob *media.Blob) error {
	panic("not implemented")
}

// DeleteBlob does not delete the file associated with it, also use DeleteFile().
func (backend *Backend) DeleteBlob(blob *media.Blob) error {
	panic("not implemented")
}

// BlobIterator returns an iterator that iterates
// all blobs as Blob structs.
func (backend *Backend) BlobIterator() model.Iterator {
	panic("not implemented")
}

// CountBlobRefs counts all BlobRef occurrences with blobID
// in all known databases.
func (backend *Backend) CountBlobRefs(blobID string) (count int, err error) {
	panic("not implemented")
}

// RemoveAllBlobRefs removes all BlobRef occurrences with blobID
// in all known databases.
func (backend *Backend) RemoveAllBlobRefs(blobID string) (count int, err error) {
	panic("not implemented")
}

// Returns ErrNotFound if no image with id is found.
func (backend *Backend) LoadImage(id string) (*media.Image, error) {
	panic("not implemented")
}

// SaveImage saves image and updates its ID if it is empty.
func (backend *Backend) SaveImage(image *media.Image) error {
	panic("not implemented")
}

func (backend *Backend) DeleteImage(image *media.Image) error {
	panic("not implemented")
}

// ImageIterator returns an iterator that iterates
// all images as Image structs.
func (backend *Backend) ImageIterator() model.Iterator {
	panic("not implemented")
}

// CountImageRefs counts all ImageRef occurrences with imageID
// in all known databases.
func (backend *Backend) CountImageRefs(imageID string) (count int, err error) {
	panic("not implemented")
}

// RemoveAllImageRefs removes all ImageRef occurrences with imageID
// in all known databases.
func (backend *Backend) RemoveAllImageRefs(imageID string) (count int, err error) {
	panic("not implemented")
}
