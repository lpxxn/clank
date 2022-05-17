package example

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/lpxxn/clank/internal"
	"gopkg.in/yaml.v2"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

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

	servSchema, err := internal.LoadSchemaFromYaml("grpc_serv.yaml")
	t.Log(servSchema)
	t.Log(servSchema.ValidateAndStartServer())
}
