package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
	"unicode"

//	"github.com/ungerik/go-start/debug"
)

func StringStartsWith(s string, start string) bool {
	return len(s) >= len(start) && s[:len(start)] == start
}

func StringEndsWith(s string, start string) bool {
	return len(s) >= len(start) && s[-len(start):] == start
}

//func FormatName(name string) string {
//	result := make([]byte, len(name))
//	first := true
//	for _, c := range name {
//		if first {
//			c
//		} else {
//		}
//	}
//	return string(result)
//}

func TextRowsAndCols(text string) (rows, cols int) {
	var x int
	for _, c := range text {
		if unicode.IsPrint(c) {
			switch c {
			case '\n':
				rows++
				x = 0
			case '\r':
				// do nothing
			default:
				x++
				if x > cols {
					cols = x
				}
			}
		}
	}
	return rows, cols
}

func StringIn(needle string, heystack []string) bool {
	for _, s := range heystack {
		if s == needle {
			return true
		}
	}
	return false
}

func ReverseStringSlice(slice []string) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func PrettifyJSON(compactJSON []byte) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, compactJSON, "", "\t")
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func MD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Todo something more generic
func StringForStruct(typeName string, attributes ...string) string {
	var builder StringBuilder

	builder.Write(typeName, "{")
	for i, attr := range attributes {
		if i > 0 {
			if i&1 == 0 {
				builder.Write(", ")
			} else {
				builder.Write(": ")
			}
		}
		builder.Write(attr)
	}
	builder.Byte('}')

	return builder.String()
}

func HasImageFileExt(filename string) bool {
	ext := path.Ext(strings.ToLower(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func IsImageURL(url string) bool {
	return HasImageFileExt(url) || StringStartsWith(url, "data:image/")
}

func JoinNonEmptyStrings(sep string, strings ...string) string {
	var buf bytes.Buffer
	first := true
	for _, s := range strings {
		if s != "" {
			if first {
				first = false
			} else {
				buf.WriteString(sep)
			}
			buf.WriteString(s)
		}
	}
	return buf.String()
}

func CompareStringsCaseInsensitive(a, b string) bool {
	return bytes.Compare(bytes.ToLower([]byte(a)), bytes.ToLower([]byte(b))) < 0
}

func EscapeJSON(jsonString string) string {
	jsonString = strings.Replace(jsonString, `\`, `\\`, -1)
	jsonString = strings.Replace(jsonString, `"`, `\"`, -1)
	return jsonString
}

func NewLineToHTML(text string) (html string) {
	return strings.Replace(text, "\n", "<br/>", -1)
}

func StripHTMLTags(text string) (plainText string) {
	chars := []byte(text)
	tagStart := -1
	for i := 0; i < len(chars); i++ {
		if chars[i] == '<' {
			tagStart = i
		} else if chars[i] == '>' && tagStart != -1 {
			chars = append(chars[:tagStart], chars[i+1:]...)
			i, tagStart = tagStart-1, -1
		}
	}
	return string(chars)
}

func ReadStringFrom(reader io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ReadBytesFrom(reader io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
