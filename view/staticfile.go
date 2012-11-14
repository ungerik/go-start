package view

import (
	"io/ioutil"
	"path"
	"time"

	"github.com/ungerik/go-start/errs"
)

func NewStaticFile(filename string) *StaticFile {
	return &StaticFile{Filename: filename}
}

// StaticFile renders a static file.
// The output is cached in memory but changes to the file on the filesystem
// cause the the cache to be rebuilt.
type StaticFile struct {
	ViewWithURLBase
	Filename string
	// Will be set automatically from Filename if empty
	ContentTypeExt string
	modifiedTime   time.Time
	cachedFileData []byte
}

func (self *StaticFile) Render(ctx *Context) (err error) {
	filePath, found, modified := FindStaticFile(self.Filename)
	if !found {
		return errs.Format("Static file not found: %s", self.Filename)
	}

	if self.ContentTypeExt == "" {
		self.ContentTypeExt = path.Ext(filePath)
	}

	if self.cachedFileData == nil || modified.After(self.modifiedTime) {
		self.cachedFileData, err = ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
	}

	ctx.Response.SetContentTypeByExt(self.ContentTypeExt)
	ctx.Response.Write(self.cachedFileData)
	return nil
}
