package model

type Configuration struct {
	Debug bool
}

var Config Configuration

const StructTagKey = "model"
