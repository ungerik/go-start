package root

import (
	"time"

	. "github.com/ungerik/go-start/view"

	"github.com/ungerik/go-start/examples/ViewPaths/views"
)

// Navigation is a function so it can be re-used by other views
func Navigation() View {
	return UL(
		// Use addresses of views variables to avoid
		// circular initialization dependencies
		A(&views.Homepage, "Home"),
		A(&views.Admin, "Admin"),
		A(&views.GetJSON, "get.json"),
		A(&views.GetXML, "get.xml"),
	)
}

func init() {
	views.Homepage = CacheView(
		time.Hour, // Cache Homepage for one hour
		&Page{
			Title: Escape("Page Title"),
			Content: Views{
				H1("go-start ViewPaths Example"),
				Navigation(),
			},
		},
	)
}
