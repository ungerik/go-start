package view

import (
	"fmt"
	"reflect"

	"github.com/ungerik/go-start/utils"
)

/*
DynamicView implements View for a function that creates and renders a dynamic
child-view in the Render method.

Example:

	dynamicView := DynamicView(
		func(response *Response) (view View, err error) {
			return HTML("return dynamic created views here"), nil
		},
	)
*/
type DynamicView func(response *Response) (view View, err error)

func (self DynamicView) Init(thisView View) {
}

func (self DynamicView) ID() string {
	return ""
}

func (self DynamicView) IterateChildren(callback IterateChildrenCallback) {
}

func (self DynamicView) Render(response *Response) error {
	child, err := self(response)
	if err != nil || child == nil {
		return err
	}
	child.Init(child)
	return child.Render(response)
}

func DynamicViewBindURLArgs(viewFunc interface{}) DynamicView {
	v := reflect.ValueOf(viewFunc)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("DynamicViewBindURLArgs: viewFunc must be a function, got %T", viewFunc))
	}
	if t.NumIn() == 0 {
		panic(fmt.Errorf("DynamicViewBindURLArgs: viewFunc has no arguments, needs at least one *view.Response"))
	}
	if t.In(0) != reflect.TypeOf((*Response)(nil)) {
		panic(fmt.Errorf("DynamicViewBindURLArgs: viewFunc's first argument must type must be *view.Response, got %s", t.In(0)))
	}
	if t.NumOut() != 2 {
		panic(fmt.Errorf("DynamicViewBindURLArgs: viewFunc must have two results, got %d", t.NumOut()))
	}
	if t.Out(0) != reflect.TypeOf((*View)(nil)).Elem() {
		panic(fmt.Errorf("RenderViewBindURLArgs: renderFunc's first result must be of type view.View, got %s", t.Out(0)))
	}
	if !utils.IsErrorType(t.Out(1)) {
		panic(fmt.Errorf("RenderViewBindURLArgs: renderFunc's second result must be of type error, got %s", t.Out(1)))
	}
	return DynamicView(
		func(response *Response) (view View, err error) {
			if len(response.Request.URLArgs) != t.NumIn()-1 {
				panic(fmt.Errorf("DynamicViewBindURLArgs: number of response URL args does not match viewFunc's arg count"))
			}
			args := make([]reflect.Value, t.NumIn())
			args[0] = reflect.ValueOf(response)
			for i, urlArg := range response.Request.URLArgs {
				val, err := utils.StringToValueOfType(urlArg, t.In(i+1))
				if err != nil {
					return nil, err
				}
				args[i+1] = reflect.ValueOf(val)
			}
			results := v.Call(args)
			if results[0].Interface() != nil {
				view = results[0].Interface().(View)
			}
			if results[1].Interface() != nil {
				err = results[1].Interface().(error)
			}
			return view, err
		},
	)
}
