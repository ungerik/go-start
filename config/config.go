package config

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/ungerik/go-start/debug"
)

var Logger = log.New(os.Stderr, "", log.LstdFlags)

type Package interface {
	Name() string
	Init() error
	io.Closer
}

var Packages []Package

func Load(configFile string, packages ...Package) {
	debug.Nop()
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		Logger.Panicf("Error while reading config file %s: %s", configFile, err)
	}
	Packages = packages
	switch path.Ext(configFile) {
	case ".json":
		// Compact JSON to make it easier to extract JSON per package
		var buf bytes.Buffer
		err = json.Compact(&buf, data)
		if err != nil {
			Logger.Panicf("Error in JSON config file %s: %s", configFile, err)
		}
		data = buf.Bytes()

		// Unmarshal packages in given order
		for _, pkg := range packages {
			// Extract JSON only for this package
			key := []byte(`"` + pkg.Name() + `":{`)
			begin := bytes.Index(data, key)
			if begin != -1 {
				begin += len(key) - 1
				end := 0
				braceCounter := 0
				for i := begin; i < len(data); i++ {
					switch data[i] {
					case '{':
						braceCounter++
					case '}':
						braceCounter--
					}
					if braceCounter == 0 {
						end = i + 1
						break
					}
				}

				err = json.Unmarshal(data[begin:end], pkg)
				if err != nil {
					Logger.Panicf("Error while unmarshalling JSON from config file %s: %s", configFile, err)
				}
			}
			err := pkg.Init()
			if err != nil {
				Logger.Panicf("Error while initializing package %s: %s", pkg.Name(), err)
			}
		}

	default:
		panic("Unsupported config file: " + configFile)
	}
}

func Close() {
	for i := len(Packages) - 1; i >= 0; i-- {
		err := Packages[i].Close()
		if err != nil {
			Logger.Println("Error while closing package %s: %s", Packages[i].Name(), err)
		}
	}
}
