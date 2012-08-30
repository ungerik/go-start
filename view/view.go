package view

type IterateChildrenCallback func(parent View, child View) (next bool)

// View is the basic interface for all types in the view package.
// A view can have an id, child views and renders its content to a XMLWriter.
// nil is permitted as View value and will be ignored while rendering HTML.
type View interface {
	Init(thisView View)
	ID() string
	IterateChildren(callback IterateChildrenCallback)
	// Everything written to out will be discarded if there was an error
	// out.Write() is not expected to return errors like bytes.Buffer 
	Render(ctx *Context) (err error)
}

// ProductionServerView returns view if view.Config.IsProductionServer
// is true, else nil which is a valid value for a View.
func ProductionServerView(view View) View {
	if !Config.IsProductionServer {
		return nil
	}
	return view
}

// NotProductionServerView returns view if view.Config.IsProductionServer
// is false, else nil which is a valid value for a View.
func NonProductionServerView(view View) View {
	if Config.IsProductionServer {
		return nil
	}
	return view
}
