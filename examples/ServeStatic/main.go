/*
## Serves static files in the current directory

* index.html will be handled as expected
* Does not list directory contents

Download, build and run example:

	go get github.com/ungerik/go-start/examples/ServeStatic
	go install github.com/ungerik/go-start/examples/ServeStatic && ServeStatic

*/
package main

import (
	"fmt"
	"path/filepath"

	"github.com/ungerik/go-start/view"
)

func main() {
	absPath, _ := filepath.Abs(".")
	fmt.Printf("Serving %s/ at http://0.0.0.0:8080/\n", absPath)
	view.Config.StaticDirs = []string{"."}
	view.RunServerAddr("0.0.0.0:8080", nil)
}
