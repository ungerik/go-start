package view

import (
	"bytes"
	"io/ioutil"
	"path"
	"time"

	"github.com/ungerik/go-start/errs"
)

func NewConcatStaticFiles(filenames ...string) *ConcatStaticFiles {
	return &ConcatStaticFiles{Filenames: filenames}
}

// ConcatStaticFiles renders multiple static files concatenated as single file.
// The output is cached in memory but changes to the files on the filesystem
// cause the the cache to be rebuilt.
type ConcatStaticFiles struct {
	ViewWithURLBase
	Filenames []string
	// Will be set automatically from Filenames[0] if empty
	ContentTypeExt  string
	modifiedTimes   []time.Time
	cachedFilesData []byte
}

func (self *ConcatStaticFiles) Render(ctx *Context) (err error) {
	if self.ContentTypeExt == "" && len(self.Filenames) > 0 {
		self.ContentTypeExt = path.Ext(self.Filenames[0])
	}
	if len(self.Filenames) != len(self.modifiedTimes) {
		self.modifiedTimes = make([]time.Time, len(self.Filenames))
	}

	for i, filename := range self.Filenames {
		_, found, modified := FindStaticFile(filename)
		if !found {
			return errs.Format("Static file not found: %s", filename)
		}
		if modified != self.modifiedTimes[i] {
			self.cachedFilesData = nil
		}
	}

	if self.cachedFilesData == nil {
		var buf bytes.Buffer
		for i, filename := range self.Filenames {
			filePath, _, modified := FindStaticFile(filename)
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			buf.Write(data)
			self.modifiedTimes[i] = modified
		}
		self.cachedFilesData = buf.Bytes()
	}

	ctx.Response.SetContentTypeByExt(self.ContentTypeExt)
	ctx.Response.Write(self.cachedFilesData)
	return nil
}
