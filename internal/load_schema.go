package internal

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadSchemaFromYaml(filePath string) (*SchemaDescription, error) {
	f, err := os.Open("grpc_serv.yaml")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	servSchema := &SchemaDescription{}
	return servSchema, yaml.Unmarshal(body, servSchema)
}
