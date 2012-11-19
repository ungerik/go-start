package view

import (
	"fmt"
	"reflect"

	"github.com/ungerik/go-start/reflection"
)

// RenderViewWithURL returns a ViewWithURL for renderFunc.
func RenderViewWithURL(renderFunc func(*Context) error) ViewWithURL {
	return NewViewURLWrapper(RenderView(renderFunc))
}

// RenderViewWithURLBindURLArgs returns a ViewWithURL for renderFunc.
// For an explanation of renderFunc see RenderViewBindURLArgs().
func RenderViewWithURLBindURLArgs(renderFunc interface{}) ViewWithURL {
	return NewViewURLWrapper(RenderViewBindURLArgs(renderFunc))
}

/*
RenderView implements all View methods for a View.Render compatible function.

Example:

	renderView := RenderView(
		func(ctx *Context) error {
			writer.Write([]byte("<html><body>Any Content</body></html>"))
			return nil
		},
	)
*/
type RenderView func(ctx *Context) error

func (self RenderView) Init(thisView View) {
}

func (self RenderView) ID() string {
	return ""
}

func (self RenderView) IterateChildren(callback IterateChildrenCallback) {
}

func (self RenderView) Render(ctx *Context) error {
	return self(ctx)
}

/*
RenderViewBindURLArgs creates a View where renderFunc is called from View.Render
with Context.URLArgs converted to function arguments.
The first renderFunc argument can be of type *Context, but this is optional.
All further arguments will be converted from the corresponding Context.URLArgs
string to their actual type.
Conversion to integer, float, string and bool arguments are supported.
If there are more URLArgs than function args, then supernumerous URLArgs will
be ignored.
If there are less URLArgs than function args, then the function args will
receive their types zero value.
The result of renderFunc can be empty or of type error.

Example:

	view.RenderViewBindURLArgs(func(i int, b bool) {
		ctx.Response.Printf("ctx.URLArgs[0]: %d, ctx.URLArgs[1]: %b", i, b)
	})
*/
func RenderViewBindURLArgs(renderFunc interface{}) RenderView {
	v := reflect.ValueOf(renderFunc)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("RenderViewBindURLArgs: renderFunc must be a function, got %T", renderFunc))
	}
	passContext := false
	if t.NumIn() > 0 {
		passContext = t.In(0) == reflect.TypeOf((*Context)(nil))
		if !passContext && !reflection.CanStringToValueOfType(t.In(0)) {
			panic(fmt.Errorf("RenderViewBindURLArgs: renderFunc's first argument must type must be of type *view.Context or convertible from string, got %s", t.In(0)))
		}
	}
	for i := 1; i < t.NumIn(); i++ {
		if !reflection.CanStringToValueOfType(t.In(i)) {
			panic(fmt.Errorf("RenderViewBindURLArgs: renderFunc's argument #%d must type must be convertible from string, got %s", i, t.In(i)))
		}
	}
	switch t.NumOut() {
	case 0:
	case 1:
		if t.Out(0) != reflection.TypeOfError {
			panic(fmt.Errorf("RenderViewBindURLArgs: renderFunc's result must be of type error, got %s", t.Out(0)))
		}
	default:
		panic(fmt.Errorf("RenderViewBindURLArgs: renderFunc must have zero or one results, got %d", t.NumOut()))
	}
	return RenderView(
		func(ctx *Context) error {
			args := make([]reflect.Value, t.NumIn())
			for i := 0; i < t.NumIn(); i++ {
				if i == 0 && passContext {
					args[0] = reflect.ValueOf(ctx)
				} else {
					iu := i
					if passContext {
						iu--
					}
					if iu < len(ctx.URLArgs) {
						val, err := reflection.StringToValueOfType(ctx.URLArgs[iu], t.In(i))
						if err != nil {
							return err
						}
						args[i] = reflect.ValueOf(val)
					} else {
						args[i] = reflect.Zero(t.In(i))
					}
				}
			}
			results := v.Call(args)
			if t.NumOut() == 1 {
				err, _ := results[0].Interface().(error)
				return err
			}
			return nil
		},
	)
}
