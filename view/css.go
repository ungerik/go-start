package view

/*
HTML5BoilerplateCSS returns a ViewWithURL that concatenates
static files with HTML5 Boilerplate CSS normalization.

Example:

	myCSS := HTML5BoilerplateCSS("common.css", "special.css")
	page := &Page{
		CSS: myCSS,
	}
*/
func HTML5BoilerplateCSS(staticCssFilenames ...string) ViewWithURL {
	staticCssFilenames = append(staticCssFilenames, "css/html5boilerplate/poststyle.css")
	staticCssFilenames = append([]string{"css/html5boilerplate/normalize.css"}, staticCssFilenames...)
	return NewConcatStaticFiles(staticCssFilenames...)
}

// NewHTML5BoilerplateCSSTemplate returns a ViewWithURL that concatenates
// text templates with HTML5 Boilerplate CSS normalization.
func NewHTML5BoilerplateCSSTemplate(getContext GetTemplateContextFunc, filenames ...string) ViewWithURL {
	views := make(Views, len(filenames)+2)
	views[0] = NewStaticFile("css/html5boilerplate/normalize.css")
	for i := range filenames {
		views[i+1] = NewTemplate(filenames[i], getContext)
	}
	views[len(views)-1] = NewStaticFile("css/html5boilerplate/poststyle.css")
	return NewViewURLWrapper(views)
}
