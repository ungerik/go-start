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
	actionsFrame := view.DIV("media-actions-frame",
		&view.Form{
			FormID:           "media-upload-image-" + name,
			SubmitButtonText: "Upload",
			// GetModel: func(form *Form, response *Response) (interface{}, error) {
			// },
		},
	)
	// &view.Button{Content: view.HTML("Change")},
	if imageRef.IsEmpty() {
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

		actionsFrame.Content = view.Views{
			actionsFrame.Content,
			&view.Form{
				FormID:           "media-remove-image-" + name,
				SubmitButtonText: "Remove",
			},
		}
	}
	return view.DIV(class,
		thumbnailFrame,
		actionsFrame,
		// description
		// link
		view.DivClearBoth(),
	), nil
}
