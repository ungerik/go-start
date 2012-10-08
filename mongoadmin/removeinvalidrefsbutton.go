package mongoadmin

import (
	"fmt"

	"github.com/ungerik/go-start/mongo"
	"github.com/ungerik/go-start/view"
)

// RemoveInvalidRefsButton returns a form with a button to remove all
// invalid mongo.Ref instances in collections.
// If len(collections) is zero, then all collections will be used.
func RemoveInvalidRefsButton(class string, collections ...*mongo.Collection) *view.Form {
	confirmation := "Are you sure you want to remove all invalid mongo refs in all collections?"
	if len(collections) > 0 {
		confirmation = "Are you sure you want to remove all invalid mongo refs in the collections "
		for i, c := range collections {
			if i > 0 {
				confirmation += ", "
			}
			confirmation += c.Name
		}
		confirmation += "?"
	}
	return &view.Form{
		FormID:              "mongoadmin.RemoveInvalidRefsButton",
		SubmitButtonText:    "Remove invalid mongo refs",
		SubmitButtonClass:   class,
		SubmitButtonConfirm: confirmation,
		OnSubmit: func(form *view.Form, formModel interface{}, ctx *view.Context) (message string, redirect view.URL, err error) {
			if len(collections) == 0 {
				for _, c := range mongo.Collections {
					collections = append(collections, c)
				}
			}
			for _, c := range collections {
				refs, err := c.RemoveInvalidRefs()
				if err != nil {
					return "", nil, err
				}
				if len(refs) > 0 {
					message = fmt.Sprintf("Removed %d invalid refs from collection %s. %s", len(refs), c.Name, message)
				}
			}
			if message == "" {
				message = "No invalid refs found"
			}
			return message, nil, nil
		},
	}
}
