package view

import (
	"github.com/ungerik/go-start/errs"
	"github.com/ungerik/go-start/utils"
	"io/ioutil"
	"path"
)

///////////////////////////////////////////////////////////////////////////////
// StaticFile

type StaticFile struct {
	ViewBase
	Filename       string // Will set file extension at ContentType
	ContentTypeExt string
	modifiedTime   int64
	fileContent    []byte
}

func (self *StaticFile) Render(context *Context, writer *utils.XMLWriter) (err error) {
	filePath, found, modified := FindStaticFile(self.Filename)
	if !found {
		return errs.Format("Static file not found: %s", self.Filename)
	}

	if self.ContentTypeExt == "" {
		self.ContentTypeExt = path.Ext(filePath)
	}

	if self.fileContent == nil || modified > self.modifiedTime {
		self.fileContent, err = ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
	}

	context.ContentType(self.ContentTypeExt)
	writer.Write(self.fileContent)
	return nil
}
