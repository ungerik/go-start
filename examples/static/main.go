/*
Build and run example:

cd examples/static/
go install && static

*/
package main

import (
	"fmt"
	"github.com/ungerik/go-start/view"
)

func main() {
	fmt.Println("Try http://0.0.0.0:8080/ -> /index.html will work as expected\n")

	// The content of multiple static directories gets merged under /
	view.Config.StaticDirs = []string{"static-dir-a", "static-dir-b"}
	view.RunServerAddr("0.0.0.0:8080", nil)
}
