package view

type Renderer interface {
	Render(ctx *Context) error
}

type Renderers []Renderer

func (self Renderers) Render(ctx *Context) error {
	for _, r := range self {
		if r != nil {
			if err := r.Render(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

type Render func(ctx *Context) error

func (self Render) Render(ctx *Context) error {
	return self(ctx)
}

// type ResponseRenderFunc func(ctx *Context) error

// func (self ResponseRenderFunc) Render(ctx *Context) error {
// 	return self(ctx)
// }

// IndirectRenderer takes the pointer to a Renderer variable
// and dereferences it when the returned Renderer's Render method is called.
// Used to break dependency cycles of variable initializations by
// using a pointer to a variable instead of its value.
func IndirectRenderer(rendererPtr *Renderer) Renderer {
	return Render(
		func(ctx *Context) (err error) {
			return (*rendererPtr).Render(ctx)
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
		func(ctx *Context) (err error) {
			if ctx.Request.Port() != port {
				return nil
			}
			return renderer.Render(ctx)
		},
	)
}

// ProductionServerRenderer returns renderer if view.Config.IsProductionServer
// is true, else nil which is a valid value for a Renderer.
func ProductionServerRenderer(renderer Renderer) Renderer {
	if !Config.IsProductionServer {
		return nil
	}
	return renderer
}

// NotProductionServerRenderer returns renderer if view.Config.IsProductionServer
// is false, else nil which is a valid value for a Renderer.
func NonProductionServerRenderer(renderer Renderer) Renderer {
	if Config.IsProductionServer {
		return nil
	}
	return renderer
}
