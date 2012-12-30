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
	// _, _, err = blobRef.TryGet()
	// if err != nil {
	// 	return nil, err
	// }

	selectModel := &view.ValueLabelSelectModel{
		ValuesAndLabels: []string{"", ""},
		SelectedValue:   blobRef.String(),
	}
	var b Blob
	for i := BlobIterator(); i.Next(&b); {
		label := fmt.Sprintf("%s (%d kb)", b.Title, b.Size.Get()/1024)
		selectModel.ValuesAndLabels = append(selectModel.ValuesAndLabels, b.ID.Get(), label)
	}
	selector := &view.Select{
		Name:  metaData.Selector(),
		Model: selectModel,
		Size:  1,
	}

	uploadList := &view.List{Class: "qq-upload-list"}

	removeButton := &view.Div{
		Class:   "qq-upload-button",
		Content: view.HTML("Remove"),
		OnClick: fmt.Sprintf(
			`jQuery("#%s").attr("value", "");
			jQuery("#%s").empty();`,
			selector.ID(),
			uploadList.ID(),
		),
	}

	uploadButton := UploadBlobButton(
		"",
		"#"+uploadList.ID(),
		fmt.Sprintf(
			`function(id, fileName, responseJSON) {
				var select = jQuery("#%s");
				select.append("<option value='"+ responseJSON.blobID +"'>" + fileName + " (" + Math.floor(responseJSON.blobSize / 1024) + " kb)</option>");
				select.attr("value", responseJSON.blobID);
			}`,
			selector.ID(),
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
		view.DIV(
			Config.BlobRefController.ActionsClass,
			selector,
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
