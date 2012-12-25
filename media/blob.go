package media

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path"

	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

func BlobIterator() model.Iterator {
	return Config.Backend.BlobIterator()
}

func NewBlob(filename string, data []byte) (*Blob, error) {
	return NewBlobFromReader(filename, bytes.NewReader(data))
}

func NewBlobFromReader(filename string, reader io.Reader) (*Blob, error) {
	contentType := mime.TypeByExtension(path.Ext(filename))
	blob := &Blob{
		Filename:    model.String(filename),
		ContentType: model.String(contentType),
	}
	blob.Init()
	writer, err := blob.FileWriter()
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	size, err := io.Copy(writer, reader)
	if err != nil {
		return nil, err
	}
	blob.Size.Set(size)
	err = blob.Save()
	if err != nil {
		return nil, err
	}
	return blob, nil
}

// NewBlobFromURL creates and saves a new Blob by downloading it from url.
func NewBlobFromURL(url string) (*Blob, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return NewBlobFromReader(path.Base(url), response.Body)
}

// LoadBlob loads an existing image from Config.Backend.
func LoadBlob(id string) (*Blob, error) {
	return Config.Backend.LoadBlob(id)
}

type Blob struct {
	ID          model.String `bson:",omitempty"`
	Title       model.String
	Link        model.Url
	Filename    model.String
	ContentType model.String
	Size        model.Int
	FileID      model.String
}

func (self *Blob) Init() *Blob {
	if self.Title == "" {
		self.Title = self.Filename
	}
	return self
}

func (self *Blob) FileURL() view.URL {
	return view.NewURLWithArgs(FileView, self.FileID.Get(), self.Filename.Get())
}

func (self *Blob) FileLink(class string) *view.Link {
	return &view.Link{
		Model: &view.URLLink{
			Url:     self.FileURL(),
			Title:   self.Title.Get(),
			Content: view.HTML(self.Filename.Get()),
		},
		Class: class,
	}
}

func (self *Blob) Save() error {
	return Config.Backend.SaveBlob(self)
}

func (self *Blob) Delete() error {
	err := self.deleteFileIfExists()
	if err != nil {
		return err
	}
	return Config.Backend.DeleteBlob(self)
}

func (self *Blob) CountRefs() (int, error) {
	return Config.Backend.CountBlobRefs(self.ID.Get())
}

func (self *Blob) RemoveAllRefs() (count int, err error) {
	return Config.Backend.RemoveAllBlobRefs(self.ID.Get())
}

// FileWriter deletes the current blob-file and returns a writer
// for a new file. The ID of the new file is set at Blob,
// but Blob is not saved. You have to call Save() after FileWriter().
func (self *Blob) FileWriter() (writer io.WriteCloser, err error) {
	err = self.deleteFileIfExists()
	if err != nil {
		return nil, err
	}
	writer, id, err := Config.Backend.FileWriter(self.Filename.Get(), self.ContentType.Get())
	if err == nil {
		self.FileID.Set(id)
	}
	return writer, err
}

// FileReader returns a reader for the blob-file if it exists.
func (self *Blob) FileReader() (reader io.ReadCloser, filename, contentType string, err error) {
	return Config.Backend.FileReader(self.FileID.Get())
}

func (self *Blob) deleteFileIfExists() error {
	if self.FileID.IsEmpty() {
		return nil
	}
	err := Config.Backend.DeleteFile(self.FileID.Get())
	if err == nil {
		self.FileID.Set("")
	}
	return err
}
