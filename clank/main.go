package main

import (
	"flag"

	"github.com/lpxxn/clank/internal"
	"github.com/lpxxn/clank/internal/clanklog"
)

var yamlPath string

func init() {
	flag.StringVar(&yamlPath, "yaml", "serv.yaml", "path to yaml file")
}

func main() {
	flag.Parse()
	servSchema, err := internal.LoadSchemaFromYaml(yamlPath)
	if err != nil {
		clanklog.Fatal(err)
	}
	if err := servSchema.ValidateAndStartServer(); err != nil {
		clanklog.Fatal(err)
	}
}
