package view

type Renderer interface {
	Render(request *Request, session *Session, response *Response) error
}

type Renderers []Renderer

func (self Renderers) Render(request *Request, session *Session, response *Response) error {
	for _, r := range self {
		if r != nil {
			if err := r.Render(request, session, response); err != nil {
				return err
			}
		}
	}
	return nil
}

type Render func(request *Request, session *Session, response *Response) error

func (self Render) Render(request *Request, session *Session, response *Response) error {
	return self(request, session, response)
}

// type ResponseRenderFunc func(response *Response) error

// func (self ResponseRenderFunc) Render(request *Request, session *Session, response *Response) error {
// 	return self(response)
// }

// IndirectRenderer takes the pointer to a Renderer variable
// and dereferences it when the returned Renderer's Render method is called.
// Used to break dependency cycles of variable initializations by
// using a pointer to a variable instead of its value.
func IndirectRenderer(rendererPtr *Renderer) Renderer {
	return Render(
		func(request *Request, session *Session, response *Response) (err error) {
			return (*rendererPtr).Render(request, session, response)
		},
	)
}

// FilterPortRenderer calls renderer.Render only
// if the request is made to a specific port
func FilterPortRenderer(port uint16, renderer Renderer) Renderer {
	if renderer == nil {
		return nil
	}
	return Render(
		func(request *Request, session *Session, response *Response) (err error) {
			if request.Port() != port {
				return nil
			}
			return renderer.Render(request, session, response)
		},
	)
}
