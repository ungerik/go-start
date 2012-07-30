package view

type CSRFProtector interface {
	ExtraFormField() View
	Validate(response *Response) (ok bool, err error)
}
