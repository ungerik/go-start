package admin

import (
	"github.com/ungerik/go-start/media"
	. "github.com/ungerik/go-start/view"

	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Admin_Images = &Page{
		Title: Escape("Images | Admin"),
		// CSS:     IndirectURL(&Admin_CSS),
		Content: Views{
			adminHeader(),
			media.ImagesAdmin(),
		},
	}
}
