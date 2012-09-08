package media

import (
	"fmt"

	"github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

func NewFormFieldFactory(wrapped view.FormFieldFactory, imageWidgetClass string, thumbnailsize int) *FormFieldFactory {
	return &FormFieldFactory{
		FormFieldFactoryWrapper: view.FormFieldFactoryWrapper{wrapped},
	}
}

type FormFieldFactory struct {
	view.FormFieldFactoryWrapper
}

func (self *FormFieldFactory) CanCreateInput(metaData *model.MetaData, form *view.Form) bool {
	debug.Nop()
	if metaData.Value.CanAddr() {
		if _, ok := metaData.Value.Addr().Interface().(*ImageRef); ok {
			return true
		}
	}
	return self.FormFieldFactoryWrapper.Wrapped.CanCreateInput(metaData, form)
}

func (self *FormFieldFactory) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (view.View, error) {
	if imageRef, ok := metaData.Value.Addr().Interface().(*ImageRef); ok {
		thumbnailSize := Config.ImageRefEditor.ThumbnailSize
		removeButton := &view.Button{Content: view.HTML("Remove"), OnClick: ""}

		var img view.View
		if imageRef.IsEmpty() {
			removeButton.Disabled = true
			img = view.IMG(Config.dummyImageURL, thumbnailSize, thumbnailSize)
		} else {
			image, err := imageRef.GetImage()
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
			Class:   Config.ImageRefEditor.ThumbnailFrameClass,
			Style:   fmt.Sprintf("width:%dpx;height:%dpx;", thumbnailSize, thumbnailSize),
			Content: img,
		}

		const onCompleteSrc = `function(id, fileName, responseJSON) {
			alert(JSON.stringify(responseJSON));
			var img = "<img src='" + responseJSON.thumbnailURL + "' width='%d' height='%d'/>";
			jQuery("#%s").empty().html(img);
			jQuery("#%s").attr("value", responseJSON.imageID);
		}`
		onComplete := fmt.Sprintf(
			onCompleteSrc,
			thumbnailSize,
			thumbnailSize,
			thumbnailFrame.ID(),
			hiddenInput.ID(),
		)
		// onComplete = ""

		editor := view.DIV(Config.ImageRefEditor.Class,
			hiddenInput,
			thumbnailFrame,
			view.DIV(Config.ImageRefEditor.ActionsClass,
				view.HTML("&larr; drag &amp; drop files here"),
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

	return self.FormFieldFactoryWrapper.Wrapped.NewInput(withLabel, metaData, form)
}
