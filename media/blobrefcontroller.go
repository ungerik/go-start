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
	blob, _, err := blobRef.TryGet()
	if err != nil {
		return nil, err
	}
	var title view.View
	if blob == nil {
		title = view.SPAN("", "Not set")
	} else {
		title = view.SPAN("", blob.Title)
	}

	hiddenInput := &view.HiddenInput{Name: metaData.Selector(), Value: blobRef.String()}

	uploadList := &view.List{Class: "qq-upload-list"}

	removeButton := &view.Div{
		Class:   "qq-upload-button",
		Content: view.HTML("Remove"),
		OnClick: fmt.Sprintf(
			`jQuery("#%s").text("Not set");
			jQuery("#%s").attr("value", "");
			jQuery("#%s").empty();`,
			title.ID(),
			hiddenInput.ID(),
			uploadList.ID(),
		),
	}

	type chooseModel struct {
		Title model.String `view:"class=gostart-select-blob"`
	}

	chooseList := view.Views{}

	uploadButton := UploadBlobButton(
		"",
		"#"+uploadList.ID(),
		fmt.Sprintf(
			`function(id, fileName, responseJSON) {
				jQuery("#%s").text(fileName);
				jQuery("#%s").attr("value", responseJSON.blobID);
			}`,
			title.ID(),
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
		view.DIV(
			Config.BlobRefController.ActionsClass,
			title,
			chooseList,
			removeButton,
			uploadButton),
		uploadList,
		view.DivClearBoth(),
	)

	if withLabel {
		return view.AddStandardLabel(form, editor, metaData), nil
	}
	return editor, nil
}
