package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var data = `
---
a: Easy!
b:
  c: 2
  d: [3, 4]
---
name: abc
age: 10
---
`

type D struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

type Student struct {
	Name string
	Age  int
}

func TestYaml(t *testing.T) {
	d2 := D{}
	decoder := yaml.NewDecoder(strings.NewReader(data))
	err := decoder.Decode(&d2)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- d2:\n%v\n\n", d2)
}

func TestYaml1(t *testing.T) {
	d := D{}
	s := Student{}
	_ = s

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

type (
	T struct {
		A string
		B string
	}
)

func Parse(source []byte) (err error) {
	dec := yaml.NewDecoder(bytes.NewReader(source))

	var doc T
	for dec.Decode(&doc) == nil {
		fmt.Println(doc)
	}

	return
}

func TestYaml3(t *testing.T) {
	s := `
---
a: val a
---
b: val b
---
`
	Parse([]byte(s))
}
