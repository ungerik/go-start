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

		hiddenInput := &view.HiddenInput{Name: metaData.Selector(), Value: imageRef.String()}

		removeButton := &view.Button{Content: view.HTML("Remove")}

		var img view.View
		if imageRef.IsEmpty() {
			removeButton.Disabled = true
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

		thumbnailFrame := &view.Div{
			Class:   Config.ImageRefEditor.ThumbnailFrameClass,
			Style:   fmt.Sprintf("width:%dpx;height:%dpx;", thumbnailSize, thumbnailSize),
			Content: img,
		}

		removeButton.OnClick = fmt.Sprintf(
			`jQuery("#%s").removeAttr("value");jQuery("#%s").empty();`,
			hiddenInput.ID(),
			thumbnailFrame.ID(),
		)

		// alert(JSON.stringify(responseJSON));
		onComplete := fmt.Sprintf(
			`function(id, fileName, responseJSON) {
				var img = "<img src='" + responseJSON.thumbnailURL + "' width='%d' height='%d'/>";
				jQuery("#%s").empty().html(img);
				jQuery("#%s").attr("value", responseJSON.imageID);
				jQuery("#%s").removeAttr("disabled");
			}`,
			thumbnailSize,
			thumbnailSize,
			thumbnailFrame.ID(),
			hiddenInput.ID(),
			removeButton.ID(),
		)

		editor := view.DIV(Config.ImageRefEditor.Class,
			hiddenInput,
			thumbnailFrame,
			view.DIV(Config.ImageRefEditor.ActionsClass,
				// view.HTML("&larr; drag &amp; drop files here"),
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
