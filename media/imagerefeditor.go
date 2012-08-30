package media

import (
	"fmt"

	"github.com/ungerik/go-start/view"
)

func ImageRefEditor(imageRef *ImageRef, name, class string, thumbnailSize int) (*view.Div, error) {
	thumbnailFrame := &view.Div{
		Class: "media-thumbnail-frame",
		Style: fmt.Sprintf("width:%dpx;height:%dpx;", thumbnailSize, thumbnailSize),
	}
	removeButton := &view.Button{Content: view.HTML("Remove")}
	actionsFrame := view.DIV("media-actions-frame",
		&view.Button{Content: view.HTML("Upload")},
		removeButton,
	)
	if imageRef.IsEmpty() {
		removeButton.Disabled = true
	} else {
		image, err := imageRef.Image()
		if err != nil {
			return nil, err
		}
		version, err := image.Thumbnail(thumbnailSize)
		if err != nil {
			return nil, err
		}
		thumbnailFrame.Content = version.ViewImage("")
	}
	return view.DIV(class,
		thumbnailFrame,
		actionsFrame,
		// description
		// link
		view.DivClearBoth(),
	), nil
}
