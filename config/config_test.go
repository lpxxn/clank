package config

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

type D struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

func TestYaml1(t *testing.T) {
	d := D{}

	err := yaml.Unmarshal([]byte(data), &d)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- d:\n%v\n\n", d)

	b, err := yaml.Marshal(&d)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- d dump:\n%s\n\n", string(b))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	b, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(b))
}

func TestYaml2(t *testing.T) {
	type configTest struct {
		Name string
		Body string
		Age  int32
	}
	type testBody struct {
		Code  string
		Data  string
		Value int32
	}
	filePath, _ := filepath.Abs("./testdata/config.yaml")
	t.Log(filePath)
	body, err := readFile(filePath)
	assert.Nil(t, err)
	c := &configTest{}
	err = yaml.Unmarshal(body, c)
	assert.Nil(t, err)
	t.Log(c)

	tb := &testBody{}
	err = json.Unmarshal([]byte(c.Body), tb)
	assert.Nil(t, err)
	t.Log(tb)
}
