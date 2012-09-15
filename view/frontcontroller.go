package view

import (
	"fmt"
	"reflect"
)

type FrontController func(content View) ViewWithURL

func (self FrontController) ContentFunc(contentFunc interface{}) ViewWithURL {
	v := reflect.ValueOf(contentFunc)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("FrontController.ForContent: contentFunc must be a function, got %s", t))
	}
	if t.NumIn() != 2 {
		panic(fmt.Errorf("FrontController.ForContent: contentFunc must have two arguments, got %d", t.NumIn()))
	}
	if t.In(0) != reflect.TypeOf((*Context)(nil)) {
		panic(fmt.Errorf("FrontController.ForContent: contentFunc's first argument must be of type *Context, got %s", t.In(0)))
	}
	if t.NumOut() != 2 {
		panic(fmt.Errorf("FrontController.ForContent: contentFunc must have two results, got %d", t.NumOut()))
	}
	if t.Out(0) != reflect.TypeOf((*View)(nil)).Elem() {
		panic(fmt.Errorf("FrontController.ForContent: contentFunc's first result must be of type View, got %s", t.Out(0)))
	}
	if t.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		panic(fmt.Errorf("FrontController.ForContent: contentFunc's second result must be of type error, got %s", t.Out(1)))
	}
	content := DynamicView(
		func(ctx *Context) (View, error) {
			if reflect.TypeOf(ctx.Data) != t.In(1) {
				panic(fmt.Errorf("FrontController: Context.Data must be of type %s, got %T", t.In(1), ctx.Data))
			}
			args := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(ctx.Data)}
			results := v.Call(args)
			view, _ := results[0].Interface().(View)
			err, _ := results[1].Interface().(error)
			return view, err
		},
	)
	return self(content)
}
