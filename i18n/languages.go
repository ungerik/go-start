package i18n

// This can be used as independent library

var iso639_1 map[string]string

func EnglishLanguageName(code string) string {
	name, ok := Languages()[code]
	if !ok {
		return code
	}
	return name
}

// ISO 639-1
func Languages() map[string]string {
	if iso639_1 == nil {
	}
	return iso639_1
}
