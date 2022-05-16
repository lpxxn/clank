package main

import (
	"flag"

	"github.com/lpxxn/clank/internal"
)

var yamlPath string

func init() {
	flag.StringVar(&yamlPath, "yaml", "serv.yaml", "path to yaml file")
}

func main() {
	flag.Parse()
	internal.Init()
	servSchema, err := internal.LoadSchemaFromYaml(yamlPath)
	if err != nil {
		panic(err)
	}
	if err := servSchema.ValidateAndStartServer(); err != nil {
		panic(err)
	}
}
