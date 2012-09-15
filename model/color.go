package model

import (
	"encoding/hex"
	"errors"
	"image/color"
	"strings"
)

func NewColor(value string) *Color {
	c := new(Color)
	c.SetString(value)
	return c
}

/*
Color holds a hex web-color with the # prefix.
The Set method accepts all valid web color syntaxes
except for color names.
Struct tag attributes:
	`model:"required"`
*/
type Color string

func (self *Color) String() string {
	return string(*self)
}

func (self *Color) SetString(s string) {
	s = strings.ToLower(s)
	switch len(s) {
	case 0, 9:
		*self = Color(s)
	case 3:
		*self = Color([]byte{'#', s[0], s[0], s[1], s[1], s[2], s[2], 'f', 'f'})
	case 4:
		if s[0] == '#' {
			*self = Color([]byte{'#', s[1], s[1], s[2], s[2], s[3], s[3], 'f', 'f'})
		} else {
			*self = Color([]byte{'#', s[0], s[0], s[1], s[1], s[2], s[2], s[3], s[3]})
		}
	case 5:
		*self = Color([]byte{'#', s[1], s[1], s[2], s[2], s[3], s[3], s[4], s[4]})
	case 6:
		*self = Color("#" + s + "ff")
	case 7:
		*self = Color(s + "ff")
	case 8:
		*self = Color("#" + s)
	default:
		panic("Invalid hex web-color length")
	}
	if self.IsValid() == false {
		panic("Invalid hex web-color")
	}
}

func (self *Color) IsEmpty() bool {
	return len(*self) == 0
}

func (self *Color) GetOrDefault(defaultColor string) string {
	if self.IsEmpty() {
		return defaultColor
	}
	return self.String()
}

func (self *Color) Required(metaData *MetaData) bool {
	return metaData.BoolAttrib(StructTagKey, "required")
}

func (self *Color) IsValid() bool {
	s := string(*self)
	l := len(s)
	if l == 0 {
		return true
	}
	if l != 9 || s[0] != '#' {
		return false
	}
	_, err := hex.DecodeString(s[1:])
	return err == nil
}

func (self *Color) Validate(metaData *MetaData) error {
	if self.Required(metaData) && self.IsEmpty() {
		return NewRequiredError(metaData)
	}
	if !self.IsValid() {
		return errors.New("Invalid hex web-color: " + string(*self))
	}
	return nil
}

// RGBA returns the color as image/color.RGBA struct.
// If the color is empty, a default zero struct will be returned.
func (self *Color) RGBA() color.RGBA {
	if self.IsEmpty() {
		return color.RGBA{}
	}
	b, err := hex.DecodeString(self.String()[1:])
	if err != nil {
		panic(err.Error())
	}
	if len(b) != 4 {
		panic("Invalid web-color length")
	}
	return color.RGBA{b[0], b[1], b[2], b[3]}
}

func (self *Color) EqualsColor(c color.Color) bool {
	r0, g0, b0, a0 := self.RGBA().RGBA()
	r1, g1, b1, a1 := c.RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1 && a0 == a1
}
