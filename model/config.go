package model

var StructTagKey = "model"

var Config Configuration

type Configuration struct {
	Debug bool
}

func (self *Configuration) Name() string {
	return "model"
}

func (self *Configuration) Init() error {
	return nil
}

func (self *Configuration) Close() error {
	return nil
}
