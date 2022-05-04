package internal

import (
	"errors"
	"fmt"
)

type SchemaDescription struct {
	Kind    ServerKind           `yaml:"kind" json:"kind"`
	Port    int                  `yaml:"port" json:"port"`
	Servers []*ServerDescription `yaml:"servers" json:"servers"`

	// gRpc
	ImportPath   []string `yaml:"importPath" json:"importPath"`
	ProtoPath    []string `yaml:"protoPath" json:"protoPath"`
	ProtosetPath string   `yaml:"protosetPath" json:"protosetPath"`
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

func (s *ServerDescription) Validate() (error, bool) {
	if s.Name == "" {
		return fmt.Errorf("server name is empty"), false
	}
	for _, m := range s.Methods {
		if err := m.Validate(); err != nil {
			return err, false
		}
	}
	return nil, true
}

type MethodDescriptionList []*MethodDescription

func (m MethodDescriptionList) Validate() error {
	for _, method := range m {
		if err := method.Validate(); err != nil {
			return err
		}
	}
}

type MethodDescription struct {
	Name            string                          `yaml:"name" json:"name"`
	DefaultResponse string                          `yaml:"defaultResponse" json:"defaultResponse"`
	Conditions      []*ResponseConditionDescription `yaml:"conditions" json:"conditions"`
}

func (m *MethodDescription) Validate() error {
	if m.Name == "" {
		return errors.New("method name is empty")
	}
	if m.DefaultResponse == "" {
		return errors.New("method default response is empty")
	}
	for _, c := range m.Conditions {
		//c.Validate()
	}
	return nil
}

type ResponseConditionDescription struct {
	Condition string `yaml:"condition" json:"condition"`
	Response  string `yaml:"response" json:"response"`
}
