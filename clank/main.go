package main

import (
	"flag"

	"github.com/lpxxn/clank/internal"
)

func main() {
	yamlPath := flag.String("yaml", "serv.yaml", "path to yaml file")
	servSchema, err := internal.LoadSchemaFromYaml(*yamlPath)
	if err != nil {
		panic(err)
	}
	if err := servSchema.ValidateAndStartServer(); err != nil {
		panic(err)
	}
}
