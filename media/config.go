package media

import (
// "github.com/ungerik/go-start/mongo"
// "github.com/ungerik/go-start/mgo"
// "errors"
)

/*
Check out:

https://github.com/valums/file-uploader
http://deepliquid.com/content/Jcrop.html
*/

var Config Configuration

type Configuration struct {
	Backend                 Backend
	NoDynamicStyleAndScript bool
}

func (self *Configuration) Name() string {
	return "media"
}

func (self *Configuration) Init() error {
	return nil
}

func (self *Configuration) Close() error {
	return nil
}
