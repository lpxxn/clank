package example

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/lpxxn/clank/internal"
	"gopkg.in/yaml.v2"
)

func TestYaml1(t *testing.T) {
	f, err := os.Open("grpc_serv.yaml")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
	m := map[string]interface{}{}
	if err := yaml.Unmarshal(body, m); err != nil {
		t.Fatal(err)
	}
	t.Log(m)

	servSchema := &internal.SchemaDescription{}
	_ = servSchema
	if err := yaml.Unmarshal(body, servSchema); err != nil {
		t.Fatal(err)
	}
	t.Log(servSchema)
	t.Log(servSchema.Validate())
	serv, err := internal.ParseServerMethodsFromProto(servSchema.ImportPath, servSchema.ProtoPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := internal.ValidateServiceInputAndOutput(servSchema.Servers, serv); err != nil {
		t.Fatal(err)
	}
	if err := internal.SetOutputFunc(servSchema.Servers, serv); err != nil {
		t.Fatal(err)
	}
	if err := serv.StartWithPort(servSchema.Port); err != nil {
		t.Fatal(err)
	}
}
