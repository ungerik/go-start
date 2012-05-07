package view

type CSRFProtector interface {
    ExtraFormField() View
    Validate(context *Context) (ok bool, err error)
}