package admin

import (
	"github.com/ungerik/go-start/media"

	. "github.com/ungerik/go-start/examples/FullTutorial/views"
)

func init() {
	Admin_Images = NewAdminPage("Images | Admin",
		media.ImagesAdmin(),
	)
}
