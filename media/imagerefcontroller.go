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
		img = view.IMG(Config.dummyImageURL, thumbnailSize, thumbnailSize)
	} else {
		image, found, err := imageRef.TryGet()
		if err != nil {
			return nil, err
		}
		if found {
			version, err := image.Thumbnail(thumbnailSize)
			if err != nil {
				return nil, err
			}
			img = version.View("")
		} else {
			imageRef.Set(nil)
			img = view.IMG(Config.dummyImageURL, thumbnailSize, thumbnailSize)
		}
	}

	hiddenInput := &view.HiddenInput{Name: metaData.Selector(), Value: imageRef.String()}

	thumbnailFrame := &view.Div{
		Class:   Config.ImageRefController.ThumbnailFrameClass,
		Style:   fmt.Sprintf("width:%dpx;height:%dpx;", thumbnailSize, thumbnailSize),
		Content: img,
	}

	uploadList := &view.List{Class: "qq-upload-list"}

	removeButton := &view.Div{
		Class:   "qq-upload-button",
		Content: view.HTML("Remove"),
		OnClick: fmt.Sprintf(
			`jQuery("#%s").attr("value", "");
			jQuery("#%s").empty().append("<img src='%s' width='%d' height='%d'/>");
			jQuery("#%s").empty();`,
			hiddenInput.ID(),
			thumbnailFrame.ID(),
			Config.dummyImageURL,
			thumbnailSize,
			thumbnailSize,
			uploadList.ID(),
		),
	}

	chooseDialogThumbnails := view.DIV("")
	chooseDialogThumbnailsID := chooseDialogThumbnails.ID()
	chooseDialog := &view.ModalDialog{
		Style: "width:600px;height:400px;",
		Content: view.Views{
			view.H3("Choose Image:"),
			chooseDialogThumbnails,
			view.ModalDialogCloseButton("Close"),
		},
	}

	chooseButton := view.DynamicView(
		func(ctx *view.Context) (view.View, error) {
			return &view.Div{
				Class:   "qq-upload-button",
				Content: view.HTML("Choose existing"),
				OnClick: fmt.Sprintf(
					`gostart_media.fillChooseDialog('#%s', '%s', function(value){
						jQuery('#%s').attr('value', value.id);
						var img = '<img src=\"'+value.url+'\" alt=\"'+value.title+'\"/>';
						jQuery('#%s').empty().append(img);
						%s
					});
					%s;`,
					chooseDialogThumbnailsID,
					API.AllThumbnails.URL(ctx.ForURLArgsConvert(Config.ImageRefController.ThumbnailSize)),
					hiddenInput.ID(),
					thumbnailFrame.ID(),
					view.ModalDialogCloseScript,
					chooseDialog.OpenScript(),
				),
			}, nil
		},
	)

	uploadButton := UploadImageButton(
		"#"+thumbnailFrame.ID(),
		"#"+uploadList.ID(),
		thumbnailSize,
		fmt.Sprintf(
			`function(id, fileName, responseJSON) {
				var img = "<img src='" + responseJSON.thumbnailURL + "' width='%d' height='%d'/>";
				jQuery("#%s").empty().html(img);
				jQuery("#%s").attr("value", responseJSON.imageID);
			}`,
			thumbnailSize,
			thumbnailSize,
			thumbnailFrame.ID(),
			hiddenInput.ID(),
		),
	)

	editor := view.DIV(Config.ImageRefController.Class,
		view.RequireScriptURL("/media/media.js", 0),
		view.RequireStyle(
			`.qq-upload-button {
				margin: 0 0 5px 10px;
				cursor: pointer;
			}
			.qq-upload-button:hover {
				background-color: #cc0000;
			}`,
			10,
		),
		chooseDialog,
		hiddenInput,
		thumbnailFrame,
		&view.Div{
			Class: Config.ImageRefController.ActionsClass,
			Style: fmt.Sprintf("margin-left: %dpx", thumbnailSize+10),
			Content: view.Views{
				removeButton,
				chooseButton,
				uploadButton,
			},
		},
		uploadList,
		view.DivClearBoth(),
	)

	if withLabel {
		return view.AddStandardLabel(form, editor, metaData), nil
	}
	return editor, nil
}
