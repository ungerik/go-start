package media

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ungerik/go-start/utils"
)

// MakePrettyUrlFilename modifies a filename so it looks good as part on an URL.
func MakePrettyUrlFilename(filename string) string {
	result := make([]byte, utf8.RuneCountInString(filename))
	i := 0
	for _, c := range filename {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.' {
			result[i] = byte(c)
		} else if c >= 'A' && c <= 'Z' {
			result[i] = byte(unicode.ToLower(c))
		} else {
			result[i] = '_'
		}
		i++
	}
	return string(result)
}

// ColoredImageDataURL creates a 1 pixel image with the given color
// and encodes it as data URL.
func ColoredImageDataURL(c color.Color) string {
	i := image.NewRGBA(image.Rect(0, 0, 1, 1))
	draw.Draw(i, i.Bounds(), image.NewUniform(c), image.ZP, draw.Src)
	var buf bytes.Buffer
	err := png.Encode(&buf, i)
	if err != nil {
		panic(err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

// ImageDataURL loads an image and encodes it as a data URL.
// Use a file-URL that begins with file:// to load local files.
func ImageDataURL(imageURL string) (dataURL string, err error) {
	var prefix string
	switch strings.ToLower(path.Ext(imageURL)) {
	case ".png":
		prefix = "data:image/png;base64,"
	case ".jpg", ".jpeg":
		prefix = "data:image/jpeg;base64,"
	case ".gif":
		prefix = "data:image/gif;base64,"
	default:
		return "", fmt.Errorf("Invalid image filename extension in URL: %s", imageURL)
	}

	data, err := utils.ReadURL(imageURL)
	if err != nil {
		return "", err
	}

	return prefix + base64.StdEncoding.EncodeToString(data), nil
}
