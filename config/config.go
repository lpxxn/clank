package config

import (
	"io/ioutil"
	"os"

	"github.com/lpxxn/clank/internal/clanklog"
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
	clanklog.Info(v, err)
	return nil
}
