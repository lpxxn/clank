package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func ReadFile(filename string) error {
	body, err := readFile(filename)
	if err != nil {
		return err
	}
	v := map[string]interface{}{}
	err = yaml.Unmarshal(body, &v)
	fmt.Println(v, err)
	return nil
}
