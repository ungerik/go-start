package media

import (
	"fmt"

	// "github.com/ungerik/go-start/debug"
	"github.com/ungerik/go-start/model"
	"github.com/ungerik/go-start/view"
)

type BlobRefController struct {
	view.SetModelValueControllerBase
}

func (self BlobRefController) Supports(metaData *model.MetaData, form *view.Form) bool {
	_, ok := metaData.Value.Addr().Interface().(*BlobRef)
	return ok
}

func (self BlobRefController) NewInput(withLabel bool, metaData *model.MetaData, form *view.Form) (input view.View, err error) {
	blobRef := metaData.Value.Addr().Interface().(*BlobRef)

	hiddenInput := &view.HiddenInput{Name: metaData.Selector(), Value: blobRef.String()}

	uploadList := &view.List{Class: "qq-upload-list"}

	removeButton := &view.Div{
		Class:   "qq-upload-button",
		Content: view.HTML("Remove"),
		OnClick: fmt.Sprintf(
			`jQuery("#%s").attr("value", "");
			jQuery("#%s").empty();`,
			hiddenInput.ID(),
			uploadList.ID(),
		),
	}

	chooseList := view.Views{
		view.JQueryUIAutocompleteFromURL(".gostart-select-blob", API.AllBlobs, 2),
	}

	uploadButton := UploadBlobButton(
		"",
		"#"+uploadList.ID(),
		fmt.Sprintf(
			`function(id, fileName, responseJSON) {
				jQuery("#%s").attr("value", responseJSON.blobID);
			}`,
			hiddenInput.ID(),
		),
	)

	editor := view.DIV(Config.BlobRefController.Class,
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
		hiddenInput,
		view.DIV(Config.BlobRefController.ActionsClass, chooseList, removeButton, uploadButton),
		uploadList,
		view.DivClearBoth(),
	)

	if withLabel {
		return view.AddStandardLabel(form, editor, metaData), nil
	}
	return editor, nil
}
