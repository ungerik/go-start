package errs

type Configuration struct {
	FormatWithCallStack bool
}

var Config = Configuration{
	FormatWithCallStack: true,
}
