package user

import (
	"github.com/ungerik/go-start/view"
)

// Nav returns user navigation links depending on if there is an user session.
// If a user session is active, logout and profile will be returned,
// if there is no active user session, login and signup will be returned.
// separator will be put between the returned links.
// signup, profile, and separator are optional and can be nil.
func Nav(login, signup, logout, profile *view.Link, separator view.View) view.View {
	return view.DynamicView(
		func(ctx *view.Context) (view.View, error) {
			if ctx.Session.ID() != "" {
				if profile == nil {
					return logout, nil
				}
				return view.Views{logout, separator, profile}, nil
			}
			if signup == nil {
				return login, nil
			}
			return view.Views{login, separator, signup}, nil
		},
	)
}
