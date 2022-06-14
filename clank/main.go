package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lpxxn/clank/internal"
	"github.com/lpxxn/clank/internal/clanklog"
)

var yamlPath string
var version = "v0.1.0"

func init() {
	flag.StringVar(&yamlPath, "yaml", "serv.yaml", "path to yaml file")
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Version: %s\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	fmt.Println(`
╭━━━┳╮╱╱╭━━━┳━╮╱╭┳╮╭━╮
┃╭━╮┃┃╱╱┃╭━╮┃┃╰╮┃┃┃┃╭╯
┃┃╱╰┫┃╱╱┃┃╱┃┃╭╮╰╯┃╰╯╯
┃┃╱╭┫┃╱╭┫╰━╯┃┃╰╮┃┃╭╮┃
┃╰━╯┃╰━╯┃╭━╮┃┃╱┃┃┃┃┃╰╮
╰━━━┻━━━┻╯╱╰┻╯╱╰━┻╯╰━╯`)
	servSchema, err := internal.LoadSchemaFromYaml(yamlPath)
	if err != nil {
		clanklog.Fatal(err)
	}
	if err := servSchema.ValidateAndStartServer(); err != nil {
		clanklog.Fatal(err)
	}
}
