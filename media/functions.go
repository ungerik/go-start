package media

import (
	"unicode"
	"unicode/utf8"

	"github.com/ungerik/go-start/view"
)

func ViewPath(name string) view.ViewPath {
	return view.ViewPath{Name: name, Sub: []view.ViewPath{
		{Name: "image", Args: 2, View: ImageView},
		{Name: "upload-image", Args: 1, View: UploadImage},
	}}
}

func MakeValidUrlFilename(filename string) string {
	result := make([]byte, utf8.RuneCountInString(filename))
	i := 0
	for _, c := range filename {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.' || c == '~' {
			result[i] = byte(c)
		} else if c >= 'A' && c <= 'Z' {
			result[i] = byte(unicode.ToLower(c))
		} else {
			result[i] = '~'
		}
		i++
	}
	return string(result)
}
