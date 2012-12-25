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

// type SubImager interface {
// 	SubImage(r image.Rectangle) image.Image
// }

func SubImageWithoutOffset(src image.Image, rect image.Rectangle) image.Image {
	switch i := src.(type) {
	case *image.RGBA:
		i = i.SubImage(rect).(*image.RGBA)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.NRGBA:
		i = i.SubImage(rect).(*image.NRGBA)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.RGBA64:
		i = i.SubImage(rect).(*image.RGBA64)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.NRGBA64:
		i = i.SubImage(rect).(*image.NRGBA64)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.YCbCr:
		i = i.SubImage(rect).(*image.YCbCr)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.Gray:
		i = i.SubImage(rect).(*image.Gray)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.Gray16:
		i = i.SubImage(rect).(*image.Gray16)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.Alpha:
		i = i.SubImage(rect).(*image.Alpha)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.Alpha16:
		i = i.SubImage(rect).(*image.Alpha16)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.Paletted:
		i = i.SubImage(rect).(*image.Paletted)
		// Remove offset
		i.Rect.Max.X = i.Rect.Dx()
		i.Rect.Min.X = 0
		i.Rect.Max.Y = i.Rect.Dy()
		i.Rect.Min.Y = 0
		return i

	case *image.Uniform:
		return i
	}
	panic(fmt.Errorf("SubImage: unsupported image type %T", src))
}

func NewImageOfType(src image.Image, width, height int) image.Image {
	return NewImageOfTypeRect(src, image.Rect(0, 0, width, height))
}

func NewImageOfTypeRect(src image.Image, bounds image.Rectangle) image.Image {
	switch i := src.(type) {
	case *image.Alpha:
		return image.NewAlpha(bounds)
	case *image.Alpha16:
		return image.NewAlpha16(bounds)
	case *image.Gray:
		return image.NewGray(bounds)
	case *image.Gray16:
		return image.NewGray16(bounds)
	case *image.NRGBA:
		return image.NewNRGBA(bounds)
	case *image.NRGBA64:
		return image.NewNRGBA64(bounds)
	case *image.Paletted:
		return image.NewPaletted(bounds, i.Palette)
	case *image.RGBA:
		return image.NewRGBA(bounds)
	case *image.RGBA64:
		return image.NewRGBA64(bounds)
	case *image.YCbCr:
		return image.NewYCbCr(bounds, i.SubsampleRatio)
	}
	panic("Unknown image type")
}
