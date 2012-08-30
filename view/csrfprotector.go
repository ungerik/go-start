package view

type CSRFProtector interface {
	ExtraFormField() View
	Validate(ctx *Context) (ok bool, err error)
}
