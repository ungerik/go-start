package root

import (
	. "github.com/ungerik/go-start/view"

	// "github.com/STARTeurope/startuplive.in/models"
	. "github.com/STARTeurope/startuplive.in/views"
)

func init() {
	Homepage = CacheView(PublicPageCacheDuration, NewPublicPage("go-start Tutorial",
		DIV("main",
			DivClearBoth(),
		),
	))
}
