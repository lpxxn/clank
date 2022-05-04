package internal

import (
	"errors"
	"fmt"
	"regexp"

	jsonIter "github.com/json-iterator/go"
)

const requestToken = "$request"

var jsonIterator = jsonIter.ConfigCompatibleWithStandardLibrary

type SchemaDescription struct {
	Kind    ServerKind            `yaml:"kind" json:"kind"`
	Port    int                   `yaml:"port" json:"port"`
	Servers ServerDescriptionList `yaml:"servers" json:"servers"`

	// gRpc
	ImportPath   []string `yaml:"importPath" json:"importPath"`
	ProtoPath    []string `yaml:"protoPath" json:"protoPath"`
	ProtosetPath string   `yaml:"protosetPath" json:"protosetPath"`
}

func (s SchemaDescription) Validate() error {
	if s.Kind != GRPC && s.Kind != HTTP {
		return errors.New("kind must be GRPC OR HTTP")
	}
	if s.Kind == GRPC {
		err := s.Servers.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type ServerKind string

const (
	GRPC ServerKind = "grpc"
	HTTP ServerKind = "http"
)

type ServerDescription struct {
	Name    string                `yaml¡:"name" json:"name"`
	Methods MethodDescriptionList `yaml:"methods" json:"methods"`
}

type ServerDescriptionList []*ServerDescription

func (s ServerDescriptionList) Validate() error {
	for _, item := range s {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerDescription) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("server name is empty")
	}
	for _, m := range s.Methods {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type MethodDescriptionList []*MethodDescription

func (m MethodDescriptionList) Validate() error {
	for _, method := range m {
		if err := method.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type MethodDescription struct {
	Name            string                          `yaml:"name" json:"name"`
	DefaultResponse string                          `yaml:"defaultResponse" json:"defaultResponse"`
	Conditions      []*ResponseConditionDescription `yaml:"conditions" json:"conditions"`
}

var re = regexp.MustCompile(`\$request.(?P<parameter>[.\w]+)`)

func (m *MethodDescription) Validate() error {
	if m.Name == "" {
		return errors.New("method name is empty")
	}
	if m.DefaultResponse == "" {
		return errors.New("method default response is empty")
	}
	params := map[string]struct{}{}
	for _, c := range m.Conditions {
		if c.Condition == "" || c.Response == "" {
			return errors.New("condition or response is empty")
		}
		match := re.FindAllStringSubmatch(c.Condition, -1)
		idx := re.SubexpIndex("parameter")
		for _, matchItem := range match {
			params[requestToken+"."+matchItem[idx]] = struct{}{}
			fmt.Println(matchItem[idx])
		}
	}
	return nil
}

type ResponseConditionDescription struct {
	Condition string `yaml:"condition" json:"condition"`
	Response  string `yaml:"response" json:"response"`
}
