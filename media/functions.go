package media

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ungerik/go-start/view"
)

// ViewPath returns the view.ViewPath for all media URLs.
func ViewPath(name string) view.ViewPath {
	return view.ViewPath{Name: name, Sub: []view.ViewPath{
		{Name: "image", Args: 2, View: ImageView},
		{Name: "upload-image", Args: 1, View: UploadImage},
	}}
}

// MakePrettyUrlFilename modifies a filename so it looks good as part on an URL.
func MakePrettyUrlFilename(filename string) string {
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
	return "data:image/png;base64," + base64.URLEncoding.EncodeToString(buf.Bytes())
}

// ImageDataURL downloads an image and encodes it as a data URL.
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

	r, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return prefix + base64.URLEncoding.EncodeToString(data), nil
}
