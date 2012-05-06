package view

type Renderer interface {
	Render(response *Response) error
}

type Renderers []Renderer

func (self Renderers) Render(response *Response) error {
	for _, r := range self {
		if r != nil {
			if err := r.Render(response); err != nil {
				return err
			}
		}
	}
	return nil
}

type Render func(response *Response) error

func (self Render) Render(response *Response) error {
	return self(response)
}

// type ResponseRenderFunc func(response *Response) error

// func (self ResponseRenderFunc) Render(response *Response) error {
// 	return self(response)
// }

// IndirectRenderer takes the pointer to a Renderer variable
// and dereferences it when the returned Renderer's Render method is called.
// Used to break dependency cycles of variable initializations by
// using a pointer to a variable instead of its value.
func IndirectRenderer(rendererPtr *Renderer) Renderer {
	return Render(
		func(response *Response) (err error) {
			return (*rendererPtr).Render(response)
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
		func(response *Response) (err error) {
			if response.Request.Port() != port {
				return nil
			}
			return renderer.Render(response)
		},
	)
}
