package view

import (
	"github.com/ungerik/go-start/utils"
	"mime/multipart"
)

type FileUploadForm struct {
	ViewBaseWithId
	Class    string
	OnSubmit func(file multipart.File, context *Context) error
}

func (self *FileUploadForm) Render(context *Context, writer *utils.XMLWriter) (err error) {
	writer.OpenTag("form").Attrib("id", self.id).AttribIfNotDefault("class", self.Class)
	writer.Attrib("action", ".").Attrib("method", "POST")
	writer.Attrib("enctype", "multipart/form-data")
	writer.OpenTag("input").Attrib("type", "file").Attrib("name", "File").CloseTag()
	writer.ForceCloseTag()
	return err
}
