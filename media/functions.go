package media

import (
	"unicode/utf8"
	"unicode"
)

func ValidUrlFilename(filename string) string {
	result := make([]byte, utf8.RuneCountInString(filename))
	i := 0
	for _, c := range filename {
		if c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '.' || c == '~' {
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
