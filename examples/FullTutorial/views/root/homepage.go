package root

import (
	"time"

	. "github.com/ungerik/go-start/view"

	// "github.com/ungerik/go-start/examples/FullTutorial/models"
	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

const PublicPageCacheDuration = time.Hour

func init() {
	Homepage = CacheView(PublicPageCacheDuration, NewPublicPage("go-start Tutorial",
		DIV("main",
			H1("Hello World from go-start"),
			DivClearBoth(),
		),
	))
}
