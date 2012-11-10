package root

import (
	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	CSS = NewHTML5BoilerplateCSSTemplate(
		TemplateContext(models.NewColorScheme()),
		"css/common.css",
		"css/style.css",
	)
}
