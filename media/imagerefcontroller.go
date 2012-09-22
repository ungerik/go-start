package media

import (
	"fmt"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type ImageRefController struct {
	view.SetModelValueControllerBase
}

func (self ImageRefController) Supports(metaData *model.MetaData, form *view.Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*ImageRef)
	return ok
}

func (self ImageRefController) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (input view.View, err error) {
	imageRef := metaData.Value.Addr().Interface().(*ImageRef)
	thumbnailSize := Config.ImageRefController.ThumbnailSize

	var img view.View
	if imageRef.IsEmpty() {
		// removeButton.Disabled = true
		img = view.IMG(Config.dummyImageURL, thumbnailSize, thumbnailSize)
	} else {
		image, err := imageRef.Get()
		if err != nil {
			return nil, err
		}
		version, err := image.Thumbnail(thumbnailSize)
		if err != nil {
			return nil, err
		}
		img = version.ViewImage("")
	}

	hiddenInput := &view.HiddenInput{Name: metaData.Selector(), Value: imageRef.String()}

	thumbnailFrame := &view.Div{
		Class:   Config.ImageRefController.ThumbnailFrameClass,
		Style:   fmt.Sprintf("width:%dpx;height:%dpx;", thumbnailSize, thumbnailSize),
		Content: img,
	}

	onComplete := fmt.Sprintf(
		`function(id, fileName, responseJSON) {
			// alert(JSON.stringify(responseJSON));
			var img = "<img src='" + responseJSON.thumbnailURL + "' width='%d' height='%d'/>";
			jQuery("#%s").empty().html(img);
			jQuery("#%s").attr("value", responseJSON.imageID);
		}`,
		thumbnailSize,
		thumbnailSize,
		thumbnailFrame.ID(),
		hiddenInput.ID(),
	)

	removeButton := &view.Button{
		Content: view.HTML("Remove"),
		OnClick: fmt.Sprintf(
			`jQuery("#%s").attr("value", ""); jQuery("#%s").empty();`,
			hiddenInput.ID(),
			thumbnailFrame.ID(),
		),
	}

	chooseDialog := &view.ModalDialog{
		Content: view.Views{
			view.H1("hello world"),
			view.ModalDialogCloseButton("Close"),
		},
	}
	// chooseDialog.Style = "width:400px;height:400px;"

	editor := view.DIV(Config.ImageRefController.Class,
		hiddenInput,
		thumbnailFrame,
		view.DIV(Config.ImageRefController.ActionsClass,
			// view.HTML("&larr; drag &amp; drop files here"),
			chooseDialog,
			chooseDialog.OpenButton("Choose"),
			removeButton,
			UploadImageButton(thumbnailSize, onComplete),
		),
		view.DivClearBoth(),
	)

	if withLabel {
		return view.AddStandardLabel(form, editor, metaData), nil
	}
	return editor, nil
}
