package models

import (
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/mongo"
)

var ColorSchemes = mongo.NewCollection("colorschemes")

func NewColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary:   "blue",
		Secondary: "black",
	}
}

type ColorScheme struct {
	mongo.DocumentBase `bson:",inline"`

	Primary   model.String
	Secondary model.String
}
