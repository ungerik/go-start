package view

import (
	"fmt"
	"reflect"

	"github.com/ungerik/go-start/reflection"
)

// DynamicViewWithURL returns a ViewWithURL for getViewFunc.
func DynamicViewWithURL(getViewFunc func(*Context) (View, error)) ViewWithURL {
	return NewViewURLWrapper(DynamicView(getViewFunc))
}

// DynamicViewWithURLBindURLArgs returns a ViewWithURL for getViewFunc.
// For an explanation of getViewFunc see DynamicViewBindURLArgs().
func DynamicViewWithURLBindURLArgs(getViewFunc interface{}) ViewWithURL {
	return NewViewURLWrapper(DynamicViewBindURLArgs(getViewFunc))
}

/*
DynamicView implements View for a function that creates and renders a dynamic
child-view in the Render method.

Example:

	dynamicView := DynamicView(
		func(ctx *Context) (view View, err error) {
			return HTML("return dynamic created views here"), nil
		},
	)
*/
type DynamicView func(ctx *Context) (view View, err error)

func (self DynamicView) Init(thisView View) {
}

func (self DynamicView) ID() string {
	return ""
}

func (self DynamicView) IterateChildren(callback IterateChildrenCallback) {
}

func (self DynamicView) Render(ctx *Context) error {
	child, err := self(ctx)
	if err != nil || child == nil {
		return err
	}
	child.Init(child)
	return child.Render(ctx)
}

/*
DynamicViewBindURLArgs creates a View where getViewFunc is called from View.Render
with Context.URLArgs converted to function arguments.
The first getViewFunc argument can be of type *Context, but this is optional.
All further arguments will be converted from the corresponding Context.URLArgs
string to their actual type.
Conversion to integer, float, string and bool arguments are supported.
If there are more URLArgs than function args, then supernumerous URLArgs will
be ignored.
If there are less URLArgs than function args, then the function args will
receive their types zero value.
The first result value of getViewFunc has to be of type View,
the second result is optional and of type error.

Example:

	view.DynamicViewBindURLArgs(func(i int, b bool) view.View {
		return view.Printf("ctx.URLArgs[0]: %d, ctx.URLArgs[1]: %b", i, b)
	})
*/
func DynamicViewBindURLArgs(getViewFunc interface{}) DynamicView {
	v := reflect.ValueOf(getViewFunc)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("DynamicViewBindURLArgs: getViewFunc must be a function, got %T", getViewFunc))
	}
	passContext := false
	if t.NumIn() > 0 {
		passContext = t.In(0) == reflect.TypeOf((*Context)(nil))
		if !passContext && !reflection.CanStringToValueOfType(t.In(0)) {
			panic(fmt.Errorf("DynamicViewBindURLArgs: getViewFunc's first argument must type must be of type *view.Context or convertible from string, got %s", t.In(0)))
		}
	}
	for i := 1; i < t.NumIn(); i++ {
		if !reflection.CanStringToValueOfType(t.In(i)) {
			panic(fmt.Errorf("DynamicViewBindURLArgs: getViewFunc's argument #%d must type must be convertible from string, got %s", i, t.In(i)))
		}
	}
	switch t.NumOut() {
	case 2:
		if t.Out(1) != reflection.TypeOfError {
			panic(fmt.Errorf("DynamicViewBindURLArgs: getViewFunc's second result must be of type error, got %s", t.Out(1)))
		}
		fallthrough
	case 1:
		if t.Out(0) != reflect.TypeOf((*View)(nil)).Elem() {
			panic(fmt.Errorf("DynamicViewBindURLArgs: getViewFunc's first result must be of type view.View, got %s", t.Out(0)))
		}
	default:
		panic(fmt.Errorf("DynamicViewBindURLArgs: getViewFunc must have one or two results, got %d", t.NumOut()))
	}
	return DynamicView(
		func(ctx *Context) (view View, err error) {
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
							return nil, err
						}
						args[i] = reflect.ValueOf(val)
					} else {
						args[i] = reflect.Zero(t.In(i))
					}
				}
			}
			results := v.Call(args)
			if t.NumOut() == 2 {
				err, _ = results[1].Interface().(error)
			}
			view, _ = results[0].Interface().(View)
			return view, err
		},
	)
}
